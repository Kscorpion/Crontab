package worker

import (
	"context"
	"github.com/Kscorpion/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//mongodb存储日志
type LogSink struct {
	client        *mongo.Client
	logCollection *mongo.Collection
	logChan       chan *common.JobLog
}

var (
	G_logSink *LogSink
)

//日志存储协程
func (logSink *LogSink) wirteLoop() {
	var (
		log *common.JobLog
	)

	for {
		select {
		case log = <-logSink.logChan:
			//把log写入MongoDB

		}
	}
}

func InitLogSink() (err error) {
	var (
		client *mongo.Client
		ctx    context.Context
	)
	//建立MongoDB连接
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(G_config.MongodbConnectTimeOut)*time.Millisecond)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(G_config.MongodbUri)); err != nil {
		return
	}

	//选择db和connection
	G_logSink = &LogSink{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
		logChan:       make(chan *common.JobLog, 1000),
	}

	//启动一个MongoDB处理协程
	go G_logSink.wirteLoop()
	return
}
