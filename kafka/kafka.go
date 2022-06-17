package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type logData struct {
	topic string
	data  string
}

//往kafka写日志的模块
var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addr []string, maxszie int) (err error) {
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
	logDataChan = make(chan *logData, maxszie)
	go sendtoKafka()//初始化就开协程等待通道传数据过来
	return
}

//向外部暴露一个函数，功能只有将日志数据发往一个通道
func SendtoChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

//私有方法，数据发往kafka的方法
func sendtoKafka() {
	for {
		select {
		case ld := <-logDataChan:
			//构造一个消息
			//sarama是操作kafka的go开发的库
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			//发送到kafka
			//不用关闭client，在后台持续调用
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send message failed, err:", err)
				return
			}
			fmt.Printf("pid:%v,offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
