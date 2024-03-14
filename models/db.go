package models

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"wozaizhao.com/wzzapi/common"
	"wozaizhao.com/wzzapi/config"
)

// DB 数据库
var DB *gorm.DB

// Models 数据库实体
var models = []interface{}{
	&User{}, &Admin{}, &Menu{}, &Role{},
}

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level Silent Info for more infomation
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  false,       // Disable color
	},
)

// DBinit 数据库初始化
func DBinit() {
	mysqlCfg := config.GetConfig().Mysql
	ds := mysqlCfg.Username + ":" + mysqlCfg.Password + "@(" + mysqlCfg.Host + ":" + mysqlCfg.Port + ")/" + mysqlCfg.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	if db, err := gorm.Open(mysql.Open(ds), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		common.LogError("DBinit", err)
		os.Exit(0)
	} else {
		DB = db
		// DB.LogMode(true)
		sqlDB, err := db.DB()
		if err != nil {
			common.LogError("db.DB", err)
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(20)
		DB.AutoMigrate(models...)
		// if err = db.AutoMigrate(models...).Error; nil != err {
		// 	config.Log("DBinit", err.Error())
		// }
	}

}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 50:
			pageSize = 50
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Search(query string, fields ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		conditions := make([]string, len(fields))
		values := make([]interface{}, len(fields))

		for i, field := range fields {
			conditions[i] = field + " LIKE ?"
			values[i] = "%" + query + "%"
		}

		return db.Where(strings.Join(conditions, " OR "), values...)
	}
}

func FieldEqual(fieldName string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" = ?", value)
	}
}

func FieldIn(fieldName string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" IN (?)", value)
	}
}
