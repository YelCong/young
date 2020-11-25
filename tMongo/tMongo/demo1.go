package tMongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "os"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`   //任务名
	Command   string    `bson:"command"`   //任务命令
	Err       string    `bson:"err"`       //错误
	Content   string    `bson:"content  "` //脚本输出
	TimePoint TimePoint `bson:"timePoint"` //执行时间点
}

type FindByJobName struct {
	JobName string `bson:"jobName"` //任务名
}

type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DelCond struct {
	beforeCond TimeBeforeCond `bson:"TimePoint.startTime"`
}

var (
	ctx context.Context
)

func TConnect() *mongo.Client {
	//mongo
	//show databases;

	var (
		err    error
		client *mongo.Client
	)
	if client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://admin:123456@localhost:27017")); err != nil {
		return nil
	}
	ctx = context.TODO()
	if err = client.Connect(ctx); err != nil {
		fmt.Println("connect fail", err)
		return nil
	}

	fmt.Println("connect success")

	return client
}

func TInsertOne() {

	var (
		collection *mongo.Collection
		client     *mongo.Client
		record     *LogRecord
		result     *mongo.InsertOneResult
		err        error
	)
	client = TConnect()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println("close fail", err)
		}
	}()

	collection = client.Database("my_db").Collection("my_collect")

	record = &LogRecord{
		JobName: "job1",
		Command: "echo 123;",
		Err:     "",
		Content: "123",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println("insert err ", err)
		return
	}
	fmt.Println("result", result.InsertedID)
}

func TFetch() {
	var (
		collection  *mongo.Collection
		database    *mongo.Database
		client      *mongo.Client
		cond        *FindByJobName
		findOptions *options.FindOptions
		cursor      *mongo.Cursor
		record      *LogRecord
		err         error
	)
	var skip int64 = 0
	var limit int64 = 2

	client = TConnect()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println("close fail", err)
		}
	}()

	database = client.Database("my_db")

	collection = database.Collection("my_collect")

	//
	cond = &FindByJobName{JobName: "job1"}
	findOptions = &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}

	if cursor, err = collection.Find(context.TODO(), cond, findOptions); err != nil {
		fmt.Println(err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		record = &LogRecord{}
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}
}

func TDelete() {
	var (
		collection *mongo.Collection
		database   *mongo.Database
		client     *mongo.Client
		err        error
		delCond    *DelCond
		delResult  *mongo.DeleteResult
	)

	client = TConnect()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println("close fail", err)
		}
	}()

	database = client.Database("my_db")
	collection = database.Collection("my_collect")
	delCond = &DelCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(delResult.DeletedCount)
}

func tCommand() {
	//mongo
	//show databases;
	//use my_db
	//
}
