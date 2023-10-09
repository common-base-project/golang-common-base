package connection

import (
	"fmt"
	"golang-common-base/pkg/config"
	"golang-common-base/pkg/logger"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

type Database struct {
	Self *gorm.DB
}

var DB *Database

func Initial() {
	DB.Init()
}

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		"Local")

	logger.Info("database url: ", config)
	//db, err := gorm.Open("mysql", config)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       config,
		DefaultStringSize:         128,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		//PrepareStmt: true,
		//Logger: logger.GetLogger(),
	})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败，连接地址: %s，error: %s", addr, err))
	}

	// 设置字符集
	//db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	// set for db connection

	sqlDB, _ := db.DB()
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(viper.GetInt(`db.gorm.maxOpenConn`))

	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
	sqlDB.SetMaxIdleConns(viper.GetInt(`db.gorm.maxIdleConn`))

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func (db *Database) Init() {
	DB = &Database{
		Self: openDB(config.DBConfig.UserName,
			config.DBConfig.Password,
			config.DBConfig.URL,
			config.DBConfig.DBName),
	}
}

func (db *Database) Close() {
	sqlDB, err := db.Self.DB()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = sqlDB.Close()
	//err := DB.Self.Close()
	if err != nil {
		logger.Error("关闭连接失败，错误信息: %s", err)
	}
}
