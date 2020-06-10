package master

import (
	"context"
	"github.com/Kscorpion/Crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

//初始化管理器
func InitLogMgr() (err error) {
	var (
		client *mongo.Client
		ctx    context.Context
	)
	//建立MongoDB连接
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(G_config.MongodbConnectTimeOut)*time.Millisecond)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(G_config.MongodbUri)); err != nil {
		return
	}

	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

func (LogMgr *LogMgr) ListLog(name string, skip int, limit int) (logArr []common.JobLog, err error) {
	var (
		filter     *common.JobLogFilter
		logSort    *common.SortLogByStartTime
		int64skip  int64
		int64limit int64
		cursor     *mongo.Cursor
		jobLog     *common.JobLog
	)

	//len(logArr)
	logArr = make([]common.JobLog, 0)

	filter = &common.JobLogFilter{JobName: name}
	int64skip = int64(skip)
	int64limit = int64(limit)
	//按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder: -1}
	//查询
	if cursor, err = LogMgr.logCollection.Find(context.TODO(), filter, &options.FindOptions{Sort: logSort}, &options.FindOptions{Skip: &int64skip}, &options.FindOptions{Limit: &int64limit}); err != nil {
		return
	}
	//延迟释放游标
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		// 反序列化BSON
		if err = cursor.Decode(jobLog); err != nil {
			continue // 有日志不合法
		}
		logArr = append(logArr, *jobLog)
	}
	//fmt.Println(logArr)
	return
}
