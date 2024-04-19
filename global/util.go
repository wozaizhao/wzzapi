package global

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
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

// MaskSensitiveInfo 对于字符串脱敏
// s 需要脱敏的字符串
// start 从第几位开始脱敏
// maskNumber 需要脱敏长度
// maskChars 掩饰字符串，替代需要脱敏处理的字符串
func MaskSensitiveInfo(s string, start int, maskNumber int, maskChars ...string) string {
	// 将字符串s的[start, end)区间用maskChar替换，并返回替换后的结果。
	maskChar := "*"
	if maskChars != nil {
		maskChar = maskChars[0]
	}
	// 处理起始位置超出边界的情况
	if start < 0 {
		start = 0
	}
	// 处理结束位置超出边界的情况
	end := start + maskNumber
	if end > len(s) {
		end = len(s)
	}
	return s[:start] + strings.Repeat(maskChar, end-start) + s[end:]
}

func GetFileNameFromUrl(url string) string {
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	return strings.Split(fileName, "?")[0]
}
