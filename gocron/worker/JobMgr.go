package worker

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
	"young/gocron/common"
)

type JobMgr struct {
	Client  *clientv3.Client
	Watcher clientv3.Watcher
	Kv      clientv3.KV
	Lease   clientv3.Lease
}

var (
	G_jobMgr JobMgr
)

func InitJobMgr() (err error) {

	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_jobMgr = JobMgr{
		Client:  client,
		Watcher: watcher,
		Lease:   lease,
		Kv:      kv,
	}

	_ = G_jobMgr.watchJob()

	//G_jobMgr.WatchKiller()

	return
}

func (jobMgr *JobMgr) watchJob() (err error) {

	var (
		watchKey           string
		watchChan          clientv3.WatchChan
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobEvent           *common.JobEvent
		jobName            string
	)

	if getResp, err = G_jobMgr.Kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return err
	}

	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			//TODO：推送给scheduler
			fmt.Println(job)
		}
	}

	go func() {
		watchKey = common.JOB_SAVE_DIR
		watchStartRevision = getResp.Header.Revision + 1

		watchChan = G_jobMgr.Watcher.Watch(context.TODO(), watchKey, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构建一个更新Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					// Delete /cron/jobs/job10
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}

					// 构建一个删除Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				//TODO：变化推送给scheduler
			}
		}
	}()

	return
}
