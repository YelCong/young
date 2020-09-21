package tEtcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"os"
	"time"
)

func TPut() {
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

	if putResp, err := kv.Put(context.TODO(), "name", "Dkngint999c"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("response", putResp.Header) //response cluster_id:14841639068965178418 member_id:10276657743932975437 revision:6 raft_term:2
	}

	if putResp, err := kv.Put(context.TODO(), "name", "yelcong", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		if putResp.PrevKv != nil {
			fmt.Println("prevValue", string(putResp.PrevKv.Value)) //prevValue Dkngint999c
		}
	}
}

func TGet() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Close()
	fmt.Println("connect success...")
	kv := clientv3.NewKV(cli)

	if response, err := kv.Get(context.TODO(), "name"); err != nil {
		fmt.Println("get err:", err)
	} else {
		fmt.Println(response.Header) //cluster_id:14841639068965178418 member_id:10276657743932975437 revision:7 raft_term:2
		fmt.Println(response.Count)  //获得的数量
		fmt.Println(response.Kvs)    //[key:"name" create_revision:2 mod_revision:7 version:6 value:"yelcong" ]
	}

	//opOption还有很多
}

func TGetDir() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Close()
	fmt.Println("connect success...")
	kv := clientv3.NewKV(cli)
	if response, err := kv.Get(context.TODO(), "/name/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Kvs)
		fmt.Printf("%T\n", response.Kvs)
		for k, v := range response.Kvs {
			fmt.Println(k, v, v.Key)
		}
	}

}

func TDelete() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Close()
	fmt.Println("connect success...")
	kv := clientv3.NewKV(cli)
	if response, err := kv.Delete(context.TODO(), "name", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.PrevKvs)
	}
}
