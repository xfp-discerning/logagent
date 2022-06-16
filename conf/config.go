package conf

import "time"

type AppConf struct {
	KafkaConf `ini:"kafka"` //这里要与配置文件对应上
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	// Topic   string `ini:"topic"`
}

type EtcdConf struct {
	Address string        `ini:"address"`
	Key     string        `ini:"collect_log_key"`
	Timeout time.Duration `ini:"timeout"`
}

//unused
type TaillogConf struct {
	Path string `ini:"path"`
}
