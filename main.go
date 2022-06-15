package main

import (
	"fmt"
	"time"

	"github.com/xfp-discerning/logagent/conf"
	"github.com/xfp-discerning/logagent/kafka"
	"github.com/xfp-discerning/logagent/taillog"
	"gopkg.in/ini.v1"
)

var cfg = new(conf.AppConf)

func run() {
	//1、读取日志
	for {
		select {
		case line := <-taillog.Readlog():
			//2、发送到kafka
			kafka.SendtoKafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}

//入口程序
func main() {
	//0、加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load ini failed, err:", err)
		return
	}
	//1、初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Println("init kafka fialed, err:", err)
		return
	}
	fmt.Println("init kafka successful")
	//2、打开日志文件准备收集日志
	err = taillog.Init(cfg.TaillogConf.Path)
	if err != nil {
		fmt.Println("init taillog failed, err:", err)
		return
	}
	fmt.Println("init taillog successful")
	run()
}
