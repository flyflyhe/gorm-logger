package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gormSession/configs"
	"gormSession/pkg/gormLogger"
	"log"
	"os"
	"sync"
	"time"
)

var gormDb *gorm.DB
var s sync.Once

func InitDb(config configs.MysqlConfig) {
	s.Do(func() {
		dsn := config.Mysql.GetDsn()
		log.Println("dsn:" + dsn)
		replicas := make([]gorm.Dialector, 0)
		for k, v := range config.MysqlReadList {
			log.Println("read", k, "dsn", v.GetDsn())
			replicas = append(replicas, mysql.Open(v.GetDsn()))
		}
		log.Println("read库", replicas)

		dbLogger := gormLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gLogger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})

		if db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{Logger: dbLogger}); err != nil {
			panic(err)
		} else {
			err := db.Use(dbresolver.Register(dbresolver.Config{
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}).SetConnMaxIdleTime(time.Hour).
				SetConnMaxLifetime(24 * time.Hour).
				SetMaxIdleConns(100).
				SetMaxOpenConns(200))
			if err != nil {
				panic(err)
			}
			gormDb = db
		}
	})
}

func GetDb() *gorm.DB {
	return gormDb
}
