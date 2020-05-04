package master

import (
	"context"
	"encoding/json"
	"github.com/Kscorpion/Crontab/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

//初始化管理器
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints, //集群地址
		DialTimeout: time.Duration(G_config.etcdDialTimeOut) * time.Millisecond,
	}

	//建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	//得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

//保存任务
func (JobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	//把任务保存到/cron/jobs/任务名->json
	var (
		jobKey      string
		jobValue    []byte
		putResponse *clientv3.PutResponse
		oldJobObj   common.Job
	)
	//etcd的保存key
	jobKey = "/cron/jobs/" + job.Name
	//任务信息json
	if jobValue, err = json.Marshal(*job); err != nil {
		return
	}
	//保存到etcd
	if putResponse, err = JobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	//如果是更新,返回旧值
	if putResponse.PrevKv != nil {
		//对旧值做一个反序列化
		if err = json.Unmarshal(putResponse.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}
