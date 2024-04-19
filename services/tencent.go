package services

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"wozaizhao.com/wzzapi/config"
	"wozaizhao.com/wzzapi/global"
)

func UploadByFile(dir string, file *multipart.FileHeader) (string, error) {
	cfg := config.GetConfig()

	// 替换为您的腾讯云 COS 配置
	secretID := cfg.TecentCosConfig.SecretId
	secretKey := cfg.TecentCosConfig.SecretKey
	region := cfg.TecentCosConfig.Region
	bucket := cfg.TecentCosConfig.Bucket

	f, err := file.Open()
	if err != nil {
		log.Errorf("file.Open Failed: %s", err)
		return "", err
	}

	defer f.Close()

	fileName := global.GetSHA1Hash(global.GenerateCode(8) + file.Filename)
	ext := filepath.Ext(file.Filename)
	key := dir + "/" + fileName + ext

	// 创建 COS 客户端
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	// 上传文件
	_, err = client.Object.Put(context.Background(), key, io.Reader(f), nil)
	if err != nil {
		log.Errorf("Object.Put Failed: %s", err)
		return "", err
	}

	return key, nil
}

func UploadByUrl(dir, fileURL string) (string, error) {
	if fileURL == "" {
		return "", fmt.Errorf("url_is_empty")
	}
	cfg := config.GetConfig()

	// 替换为您的腾讯云 COS 配置
	secretID := cfg.TecentCosConfig.SecretId
	secretKey := cfg.TecentCosConfig.SecretKey
	region := cfg.TecentCosConfig.Region
	bucket := cfg.TecentCosConfig.Bucket

	key := dir + "/" + global.GetFileNameFromUrl(fileURL)
	// 获取文件内容
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Errorf("http.Get Failed: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("io.ReadAll Failed: %s", err)
		return "", err
	}

	// 创建 COS 客户端
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	// 上传文件
	_, err = client.Object.Put(context.Background(), key, bytes.NewReader(fileContent), nil)
	if err != nil {
		log.Errorf("Object.Put Failed: %s", err)
		return "", err
	}

	return key, nil
}
