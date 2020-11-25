package tEtcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"os"
	"time"
)

func TSetLease() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()

	//申请一个lease
	lease := clientv3.NewLease(cli)

	//申请一个10s的租约
	leaseGrantResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	leaseID := leaseGrantResp.ID

	//put一个kv，让它和kv关联起来，从而实现有效期
	kv := clientv3.NewKV(cli)
	if putResp, err := kv.Put(context.TODO(), "/cron/lock/job1", "1", clientv3.WithLease(leaseID)); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(putResp)
	}

	//定时看是否过期
	for {
		if getResp, err := kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		} else {
			if getResp.Count == 0 {
				fmt.Println("expired...")
				break
			}
			fmt.Println("有效期内", getResp.Kvs)
			time.Sleep(2 * time.Second)
		}
	}
}

func TWatch() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 1 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(101)
	}
	defer cli.Close()
	fmt.Println("connect success...")

	kv := clientv3.NewKV(cli)

	//模拟etcd中kv的变化
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job 7 ")
			time.Sleep(time.Second)
			kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(time.Second)
		}
	}()
	getResp, err := kv.Get(context.TODO(), "/cron/jobs/job7")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值", string(getResp.Kvs[0].Value))
	}

	//当前etcd事务编号，单调递增，从这个版本开始监听
	watchStartRevision := getResp.Header.Revision + 1

	//创建监听器
	watcher := clientv3.Watcher(cli)

	//创建一个上下文来终止监听(5秒后自动终止)
	ctx,cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5 * time.Second, func() {
		cancelFunc()
	})

	watchRespChan := watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为", string(event.Kv.Value), event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了",event.Kv.ModRevision)
			}
		}
	}

}
