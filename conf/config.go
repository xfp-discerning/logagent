package conf

type AppConf struct {
	KafkaConf   `ini:"kafka"` //这里要与配置文件对应上
	TaillogConf `ini:"taillog"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type TaillogConf struct {
	Path string `ini:"path"`
}
