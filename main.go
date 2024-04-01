package main

import (
	log "github.com/sirupsen/logrus"
	"time"

	"wozaizhao.com/wzzapi/config"
	"wozaizhao.com/wzzapi/global"
	"wozaizhao.com/wzzapi/models"
	"wozaizhao.com/wzzapi/router"
)

func main() {

	r := router.SetupRouter()
	r.SetTrustedProxies([]string{"0.0.0.0/0", "127.0.0.1"})
	cfg := config.GetConfig()

	if cfg.Mode == "production" {
		global.LogToFile()
	}
	// 加载东八区时区
	sh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Errorf("LoadLocation Failed: %s", err)
		return
	} else {
		global.SetTimeZone(sh)
	}

	models.SetKey([]byte(cfg.EncryptionKey))
	models.DBinit()

	r.Run(cfg.Listen)
}
