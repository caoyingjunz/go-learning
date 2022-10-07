package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer etcd.Close()

	// 写入 etcd 的值
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = etcd.Put(ctx, "/config/name", "caoyingjun")
	if err != nil {
		panic(err)
	}

	// 查询 etcd 的值
	resp, err := etcd.Get(ctx, "/config/name")
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// etcd支持key前缀匹配，添加 clientv3.WithPrefix()参数即可
	// 读取key前缀等于"/config"的所有值
	resp, err = etcd.Get(ctx, "/config", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	// 遍历查询结果
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	// 删除操作
	_, err = etcd.Delete(ctx, "/config/name")
	// 前缀删除
	_, err = etcd.Delete(ctx, "/config", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	// 监听效果
	// 测试写入
	go func() {
		for {
			_, _ = etcd.Put(context.Background(), "/config/name", time.Now().String())
			time.Sleep(2 * time.Second)
		}
	}()
	wChan := etcd.Watch(context.Background(), "/config/name", clientv3.WithPrefix()) // 监听key前缀的一组key的值
	for watchResp := range wChan {
		for _, event := range watchResp.Events {
			fmt.Printf("Event received! %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}
}
