package main

import (
	"time"

	"wozaizhao.com/wzzapi/common"
	"wozaizhao.com/wzzapi/config"
	"wozaizhao.com/wzzapi/models"
	"wozaizhao.com/wzzapi/router"
)

func main() {

	r := router.SetupRouter()
	r.SetTrustedProxies([]string{"0.0.0.0/0", "127.0.0.1"})
	cfg := config.GetConfig()

	if cfg.Mode == "production" {
		common.LogToFile()
	}
	// 加载东八区时区
	sh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		common.LogError("LoadLocation", err)
		return
	} else {
		common.SetTimeZone(sh)
	}

	models.SetKey([]byte(cfg.EncryptionKey))
	models.DBinit()

	r.Run(cfg.Listen)
}
