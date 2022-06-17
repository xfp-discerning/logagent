package taillog

import "github.com/xfp-discerning/logagent/etcd"

type taillogmgr struct {
	logEntryConf []*etcd.LogEntry
	// taskMap map[string] *TailTask
}

var tailmgr *taillogmgr

func Init(logEntryConf []*etcd.LogEntry){
	tailmgr = &taillogmgr{
		logEntryConf: logEntryConf,
	}
	for _, logetry := range  tailmgr.logEntryConf {//视频中的tailmgr.logEntryConf为logEntryConf
		//logEntry.path是要收集日志的路径
		NewTailTask(logetry.Path, logetry.Topic)
	}
}
