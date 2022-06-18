package taillog

import (
	"fmt"
	"time"

	"github.com/xfp-discerning/logagent/etcd"
)

type taillogmgr struct {
	logEntryConf []*etcd.LogEntry
	taskMap      map[string]*TailTask
	newConfChan  chan []*etcd.LogEntry
}

var tailmgr *taillogmgr

func Init(logEntryConf []*etcd.LogEntry) {
	tailmgr = &taillogmgr{
		logEntryConf: logEntryConf, //将配置文件收集起来保存起来
		taskMap:      make(map[string]*TailTask, 16),
		newConfChan:  make(chan []*etcd.LogEntry),
	}
	for _, logentry := range logEntryConf { //视频中的tailmgr.logEntryConf为logEntryConf
		//logEntry.path：要收集日志的路径
		//初始化时，启了多少个tailtask都要记下来，为了后续判断方便
		tailObj := NewTailTask(logentry.Path, logentry.Topic)
		mk := fmt.Sprintf("%s_%s", logentry.Path, logentry.Topic)
		tailmgr.taskMap[mk] = tailObj
	}
	go tailmgr.run()
}

//监听自己的newConfChan，有了新的配置过来就做对应的处理
func (t *taillogmgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			for _, conf := range newConf {
				//1、配置更新
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.taskMap[mk]
				if ok {
					//原来就有，不要操作
					continue
				} else {
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.taskMap[mk] = tailObj
				}
			}
			//找出原来t.logEntry有，但是newconf没有的，要删除
			for _, c1 := range t.logEntryConf { //从原来的配置中依次拿出配置项
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					//把c1对应的这个tailObj给停掉
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.taskMap[mk].cancelFunc()
				}
			} 
			//2、配置删除
			//3、配置变更
			fmt.Println("新的配置来了", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

func PushNewConf() chan<- []*etcd.LogEntry {
	return tailmgr.newConfChan
}
