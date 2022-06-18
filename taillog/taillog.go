package taillog

import (
	"context"
	"fmt"

	"github.com/hpcloud/tail"
	"github.com/xfp-discerning/logagent/kafka"
)

//从日志文件收集日志的模块

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	//为了实现退出t.run
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path string, topic string) (tailObj *TailTask) {
	ctx,cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:  path,
		topic: topic,
		ctx: ctx,
		cancelFunc: cancel,
	}
	tailObj.init() //according to path to open log//利用tailf包
	return
}

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed ,err: ", err)
		return
	}
	go t.run()
}

//3.2发往kafka
func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tailtask : %s_%s has exited\n",t.path,t.topic)
			return
		case line := <-t.instance.Lines:
			//发往kafka
			// kafka.SendtoKafka(t.topic, line.Text)//函数调用函数
			//先把日志数据发往一个通道中
			kafka.SendtoChan(t.topic, line.Text)
			//kafka中有单独的协程去通道取日志
		}
	}
}
