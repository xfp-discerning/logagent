package taillog

import (
	"fmt"

	"github.com/hpcloud/tail"
)

//从日志文件收集日志的模块

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path string, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}
	tailObj.Init() //according to path to open log
	return
}

func (t *TailTask) Init() {
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
}

func (t *TailTask) Readlog() <-chan *tail.Line {
	return t.instance.Lines
}
