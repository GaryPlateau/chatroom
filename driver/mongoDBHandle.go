package driver

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"michatroom/conf"
	"michatroom/utils"
)

var MongoDBSingleInstance *MongoDBHandler

type MongoDBHandler struct {
	mongoDBClient *mongo.Client
}

func NewMongoDBInstance() *MongoDBHandler {
	if MongoDBSingleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		MongoDBSingleInstance = new(MongoDBHandler)
		MongoDBSingleInstance.initMongoDB()
	}
	return MongoDBSingleInstance
}

func (this *MongoDBHandler) CloseMongoDB() {
	this.mongoDBClient.Disconnect(context.TODO())
}

func (this *MongoDBHandler) initMongoDB() {
	var err error
	// 设置mongoDB客户端连接信息
	//mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	clientOptions := options.Client().ApplyURI(conf.MongoDBDSN)
	//option := options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetConnectTimeout(2 *time.Second).SetAuth(options.Credential{Username:"tester",Password:"123",AuthSource: "test"})
	this.mongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		return
	}

	err = this.mongoDBClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("MongoDB Connect")
	return
}

func (this *MongoDBHandler) CraeteCollection(dbName string, collectionName string) *mongo.Collection {
	return this.mongoDBClient.Database(dbName).Collection(collectionName)
}

func (this *MongoDBHandler) InsertDb(ctx context.Context, collection *mongo.Collection, doc map[string]interface{}) {
	res, err := collection.InsertOne(ctx, doc)
	utils.ErrorHandler("新增错误:", err)
	fmt.Printf("insert id %v \n", res.InsertedID)

	//docs := []interface{}{Student{Name: "李四", Age: 15, Score: 50}, Food{Sweet: 70.1, Spices: 0.0, Salty: 10.3}}
	//ress, err := collection.InsertMany(ctx, docs)
	//fmt.Printf("insert many ids:%v\n", ress.InsertedIDs)
}

func (this *MongoDBHandler) Update(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{"city", "北京"}}
	update := bson.D{{"$inc", bson.D{{"score", 5}}}} //inc为increase
	res, err := collection.UpdateMany(ctx, filter, update)
	utils.ErrorHandler("更新错误:", err)
	fmt.Printf("update %d doc \n", res.ModifiedCount)
}

func (this *MongoDBHandler) Query(ctx context.Context, collection *mongo.Collection) {
	sort := bson.D{{"name", 1}} //1为升序
	filter := bson.D{{"score", bson.D{{"$gt", 3}}}}
	findOption := options.Find()
	findOption.SetSort(sort)
	findOption.SetLimit(10)
	findOption.SetSkip(1)
	_, err := collection.Find(ctx, filter, findOption)
	utils.ErrorHandler("查询错误:", err)
	//for cursor.Next(ctx) {
	//	var doc Student
	//	err := cursor.Decode(&doc)
	//	CheckError("查询错误:", err)
	//	fmt.Printf("%s %d %d\n", doc.Name, doc.Age, doc.Score)
	//}
}

func (this *MongoDBHandler) Delete(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{"name", "张三"}}
	res, err := collection.DeleteMany(ctx, filter)
	utils.ErrorHandler("删除错误:", err)
	fmt.Printf("delete %d doc \n", res.DeletedCount)
}
