package common

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"time"
)

var TimeZone *time.Location

func GenerateCode(l int) string {
	// 生成一个 32 字节的随机字节数组
	bytes := make([]byte, l)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	// 将字节数组转换为 base64 编码的字符串
	code := base64.URLEncoding.EncodeToString(bytes)

	return code
}

func GetSHA1Hash(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func ParseInt(str string) (int, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	return int(num), err
}

func SetTimeZone(timeZone *time.Location) {
	TimeZone = timeZone
}
