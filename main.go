package main

import (
	"fmt"
	"time"

	"github.com/xfp-discerning/logagent/conf"
	"github.com/xfp-discerning/logagent/etcd"
	"github.com/xfp-discerning/logagent/kafka"
	"github.com/xfp-discerning/logagent/taillog"
	"gopkg.in/ini.v1"
)

var cfg = new(conf.AppConf)

// func run() {
// 	//1、读取日志
// 	for {
// 		select {
// 		case line := <-taillog.Readlog():
// 			//2、发送到kafka
// 			kafka.SendtoKafka(cfg.KafkaConf.Topic, line.Text)
// 		default:
// 			time.Sleep(time.Second)
// 		}
// 	}
// }

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
	//2、初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd failed,err:", err)
		return
	}
	fmt.Println("init etcd success")

	//2.1从etcd中获取日志收集项的配置信息
	//2.2用哨兵去监视日志收集项的变化（有变化即使通知logagent实现热加载）
	logEntryConf, err := etcd.Getconf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Println("get conf from etcd failed, err:", err)
		return
	}
	fmt.Printf("logEntryConf:%v\n", logEntryConf)
	for index, value := range logEntryConf {
		fmt.Printf("index:%v, value:%v\n", index, value)
	}
	//3、收集日志发往kafka
	//3.1循环每一个日志项，创建tailObj
	for _, logetry := range logEntryConf {
		//logEntry.path是要收集日志的路劲
		task := taillog.NewTailTask(logetry.Path, logetry.Topic)
		for {
			select {
			case line := <-task.Readlog():
				kafka.SendtoKafka(logetry.Topic, line.Text)
			}
		}
	}
	//3.2发往kafka

	// //2、打开日志文件准备收集日志
	// err = taillog.Init(cfg.TaillogConf.Path)
	// if err != nil {
	// 	fmt.Println("init taillog failed, err:", err)
	// 	return
	// }
	// fmt.Println("init taillog successful")
	// run()
}
