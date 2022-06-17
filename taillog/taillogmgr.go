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
	for _, logentry := range tailmgr.logEntryConf { //视频中的tailmgr.logEntryConf为logEntryConf
		//logEntry.path是要收集日志的路径
		NewTailTask(logentry.Path, logentry.Topic)
	}
	go tailmgr.run()
}

//监听自己的newConfChan，有了新的配置过来就做对应的处理
func (t *taillogmgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			//1、配置更新
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
