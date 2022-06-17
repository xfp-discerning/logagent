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

var cfg = new(conf.AppConf) //加载配置文件的全局变量

//入口程序
func main() {
	//0、加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load ini failed, err:", err)
		return
	}
	//1、初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address},cfg.KafkaConf.Chan_max_size)
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
	taillog.Init(logEntryConf)
}
