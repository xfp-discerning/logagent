package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type LogEntry struct {
	Path  string `json:"path"` //tag，可以用json.Unmarshal直接反序列化存到对象中
	Topic string `json:"topic"`
}

var cli *clientv3.Client

func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("connect etcd failed, err:", err)
		return
	}
	return
}

func Getconf(key string) (logEntryConf []*LogEntry, err error) {
	//get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("get from etcd failed, err:", err)
		return
	}
	for _, kv := range resp.Kvs {
		// fmt.Printf("%s : %s\n", kv.Key, kv.Value)
		err = json.Unmarshal(kv.Value, &logEntryConf)
		if err != nil {
			fmt.Println("unmarshal etcd value failed, err:", err)
			return
		}
	}
	return
}

//哨兵
func WatchConf(key string, NewConfCh chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), key) //保持监控不断开
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("type:%v, key:%v, value:%v\n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			var NewConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				//如果是删除操作，手动传递一个空的配置项
				err := json.Unmarshal(evt.Kv.Value, &NewConf)
				if err != nil {
					fmt.Println("unmarshal failed, err:", err)
					continue
				}
			}
			fmt.Printf("get new config: %v\n", NewConf)
			NewConfCh <- NewConf
		}
	}
}
