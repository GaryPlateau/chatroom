package main

import (
	"michatroom/conf"
	"michatroom/driver"
	"michatroom/router"
)

func init() {
	conf.ConfigInit()
	driver.MysqlSingleInstance = driver.NewMysqlInstance()
	driver.MongoDBSingleInstance = driver.NewMongoDBInstance()
	driver.RedisSingleInstance = driver.NewRedisInstance()
}

func destory() {
	driver.MysqlSingleInstance.Close()
	driver.MongoDBSingleInstance.CloseMongoDB()
	driver.RedisSingleInstance.Close()
}

func main() {
	router.ManRouter()
	defer destory()
}
