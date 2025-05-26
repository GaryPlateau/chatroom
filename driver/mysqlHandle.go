package driver

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"michatroom/conf"
	"os"
	"os/signal"
	"time"
)

var MysqlSingleInstance *MysqlHandler

type MysqlHandler struct {
	Db *gorm.DB
}

func NewMysqlInstance() *MysqlHandler {
	if MysqlSingleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		MysqlSingleInstance = new(MysqlHandler)
		MysqlSingleInstance.initMysql()
	}
	return MysqlSingleInstance
}

func (this *MysqlHandler) Close() {
	sqlDB, _ := this.Db.DB()
	sqlDB.Close()
}

func (this *MysqlHandler) initMysql() {
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql日志
			LogLevel:      logger.Info, //日志级别
			Colorful:      true,        //日志颜色
		},
	)
	this.Db, err = gorm.Open(mysql.Open(conf.MysqlDSN), &gorm.Config{
		//DryRun: false, 直接运行生成sql而不执行
		Logger: newLogger, // Logger: logger.Default.LogMode(logger.Info), // 可以打印SQL
		// QueryFields: true, // 使用表的所有字段执行SQL查询
		// 关闭事务，gorm默认是打开的，关闭可以提升性能
		// SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("打开数据库链接失败:", err)
		return
	}

	c := make(chan os.Signal)
	if err != nil {
		fmt.Println("Connect database failed and Shutdown web server")
		c <- os.Interrupt
		signal.Notify(c, os.Interrupt)
		// 下面是监听所有的信号
		// signal.Notify(c)
		// 致命错误，如果web服务器在它启动前，则不能关闭服务器，所以需要用到上面的signal 来处理
		// log.Fatal("connect database failed", err)
		return
	}

	sqlDB, _ := this.Db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	return
}
