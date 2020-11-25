package tEtcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"os"
	"time"
)

func TOp() {
	var err error

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(101)
	}
	defer cli.Close()
	fmt.Println("connect success...")

	kv := clientv3.NewKV(cli)

	//创建op：operation
	putOp := clientv3.OpPut("/cron/jobs/job8", "11")

	//执行op
	var (
		opResp clientv3.OpResponse
	)

	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(opResp.Put().Header.Revision)

	getOp := clientv3.OpGet("/cron/jobs/job8")

	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("数据Revision", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据Value", string(opResp.Get().Kvs[0].Value))
}
