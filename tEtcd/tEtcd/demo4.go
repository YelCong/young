package tEtcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func DistributedOptimisticLock() {
	fmt.Println("分布式乐观锁...")
	fmt.Println("lease实现锁的自动过期")
	fmt.Println("op+txn事务")

	var (
		err                        error
		config                     clientv3.Config
		client                     *clientv3.Client
		lease                      clientv3.Lease
		leaseGrantResponse         *clientv3.LeaseGrantResponse
		leaseID                    clientv3.LeaseID
		chanLeaseKeepAliveResponse <-chan *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveResponse     *clientv3.LeaseKeepAliveResponse
		ctx                        context.Context
		cancelFunc                 context.CancelFunc
		kv                         clientv3.KV
		txn                        clientv3.Txn
		txnResp                    *clientv3.TxnResponse
	)
	//0.建立连接
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	client, err = clientv3.New(config)
	if err != nil {
		fmt.Println("connect fail...", err)
		return
	}
	defer client.Close()
	fmt.Println("connect success...")

	//1.上锁(创建租约，自动续租，拿着租约去抢占一个key)
	lease = clientv3.NewLease(client)
	if leaseGrantResponse, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	leaseID = leaseGrantResponse.ID

	ctx, cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc() //确保函数退出后，停止自动续租
	defer func() {
		if _, err := lease.Revoke(context.TODO(), leaseID); err != nil {
			fmt.Println(err)
			return
		}
	}()

	if chanLeaseKeepAliveResponse, err = lease.KeepAlive(ctx, leaseID); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case leaseKeepAliveResponse = <-chanLeaseKeepAliveResponse:
				if leaseKeepAliveResponse == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {
					fmt.Println("收到自动续租应答：", leaseKeepAliveResponse.ID)
				}
			}
		}
	END:
	}()

	//抢锁

	kv = clientv3.NewKV(client)
	txn = kv.Txn(context.TODO())

	//事务
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet("/cron/lock/job9"))

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	//判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}

	//2.处理业务
	fmt.Println("doing job")
	time.Sleep(5 * time.Second)

	//3.通过defer释放资源
}
