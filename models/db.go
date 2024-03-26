package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	_log "log"
	"os"
	"strings"
	"time"
	"wozaizhao.com/wzzapi/config"
)

// DB 数据库
var DB *gorm.DB

// 密钥
var key []byte

func SetKey(val []byte) {
	key = val
}

// Models 数据库实体
var models = []interface{}{
	&User{}, &Admin{}, &Menu{}, &Role{},
}

var newLogger = logger.New(
	_log.New(os.Stdout, "\r\n", _log.LstdFlags), // io writer
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
		log.Errorf("DBinit Failed: %s", err)
		os.Exit(0)
	} else {
		DB = db
		// DB.LogMode(true)
		sqlDB, err := db.DB()
		if err != nil {
			log.Errorf("DBinit Failed: %s", err)
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

func encrypt(message string) (encoded string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)

	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func decrypt(secure string) (decoded string, err error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext_block_size_is_too_short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
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
