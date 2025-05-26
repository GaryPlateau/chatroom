package conf

import (
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
)

var (
	MysqlDSN   string
	RedisOpt   *redis.Options
	MongoDBDSN string
)

var (
	HttpAddr      string
	HttpPort      string
	MaxExpireTime int

	db         string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	dbCharset  string

	redisDb          string
	redisAddr        string
	redisPassword    string
	redisDbName      int
	redisPoolSize    int
	redisMinIdleConn int

	mongoDBName string
	mongoDBAddr string
	mongoDBPwd  string
	mongoDBPort string
)

func ConfigInit() {
	//从本地读取环境变量
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadMysqlData(file)
	LoadRedisData(file)
	LoadMongoDB(file)

	//MySQL
	//manager:Package main@tcp(162.14.67.106:3306)/chatroom?charset=utf8mb4&parseTime=True&loc=Local
	MysqlDSN = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName, dbCharset)
	//path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	//Redis
	RedisOpt = &redis.Options{
		Addr:         redisAddr, // url
		Password:     redisPassword,
		DB:           redisDbName, // 0号数据库
		MinIdleConns: redisMinIdleConn,
		PoolSize:     redisPoolSize,
	}
	//MongoDB
	//mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	MongoDBDSN = fmt.Sprintf("mongodb://%v:%v", mongoDBAddr, mongoDBPort)
}

func LoadServer(file *ini.File) {
	HttpAddr = file.Section("service").Key("HttpAddr").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
	MaxExpireTime, _ = file.Section("service").Key("MaxExpireTime").Int()
}

func LoadMysqlData(file *ini.File) {
	db = file.Section("mysql").Key("Db").String()
	dbHost = file.Section("mysql").Key("DbHost").String()
	dbPort = file.Section("mysql").Key("DbPort").String()
	dbUser = file.Section("mysql").Key("DbUser").String()
	dbPassword = file.Section("mysql").Key("DbPassword").String()
	dbName = file.Section("mysql").Key("DbName").String()
	dbCharset = file.Section("mysql").Key("DbCharset").String()
}

func LoadRedisData(file *ini.File) {
	redisDb = file.Section("redis").Key("RedisDb").String()
	redisAddr = file.Section("redis").Key("RedisAddr").String()
	redisPassword = file.Section("redis").Key("RedisPassword").String()
	redisDbName, _ = file.Section("redis").Key("RedisDbName").Int()
	redisPoolSize, _ = file.Section("redis").Key("RedisPoolSize").Int()
	redisMinIdleConn, _ = file.Section("redis").Key("RedisMinIdleConn").Int()
}

func LoadMongoDB(file *ini.File) {
	mongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	mongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	mongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	mongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
}
