package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

//往kafka写日志的模块
var client sarama.SyncProducer

func Init(addr []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //发送数据需要leader和follower都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出来一个partition
	config.Producer.Return.Successes = true                   //成功交付信息将在success channel返回

	//连接kafka
	client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		fmt.Println("producer closed,err:", err)
		return
	}
	return
}

func SendtoKafka(topic,data string){
	//构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	//发送到kafka
	//不用关闭client，在后台持续调用
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed, err:", err)
		return
	}
	fmt.Printf("pid:%v,offset:%v\n", pid, offset)
}
