package driver

import (
	"fmt"
	"github.com/go-redis/redis"
	"michatroom/conf"
	"michatroom/utils"
)

var RedisSingleInstance *RedisHandler

type RedisHandler struct {
	redisClient *redis.Client
}

func NewRedisInstance() *RedisHandler {
	if RedisSingleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		RedisSingleInstance = new(RedisHandler)
		RedisSingleInstance.initRedis()
	}
	return RedisSingleInstance
}

func (this *RedisHandler) Close() {
	this.redisClient.Close()
}

func (this *RedisHandler) initRedis() {
	this.redisClient = redis.NewClient(conf.RedisOpt)
	result, err := this.redisClient.Ping().Result()
	utils.ErrorHandler("ping err:", err)
	fmt.Println("initRedis:", result)
}

func (this *RedisHandler) SetValue(key string, value string) {

	err := this.redisClient.Set(key, value, 0).Err()
	utils.ErrorHandler("setvalue is err:", err)
}

func (this *RedisHandler) GetValue(key string) string {
	result, err := this.redisClient.Get(key).Result()
	utils.ErrorHandler("getvalue is err:", err)
	return result
}

func (this *RedisHandler) HSetValue(key, field string, value interface{}) {
	err := this.redisClient.HSet(key, field, value).Err()
	utils.ErrorHandler("HSetValue is err :", err)
}

func (this *RedisHandler) HGetValue(key string, field string) interface{} {
	result, err := this.redisClient.HGet(key, field).Result()
	utils.ErrorHandler("HGetValue is err:", err)
	return result
}

func (this *RedisHandler) HMSetValue(key string, values map[string]interface{}) {
	err := this.redisClient.HMSet(key, values).Err()
	utils.ErrorHandler("HMSetValue is err :", err)
}
