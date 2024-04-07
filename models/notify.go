package models

import (
	// log "github.com/sirupsen/logrus"
	"time"

	"gorm.io/gorm"
)

// 通知
type Notify struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Sender      uint           `json:"sender" gorm:"type:varchar(64);DEFAULT '';comment:发送通道ID"`
	Title       string         `json:"title" gorm:"type:varchar(255);DEFAULT '';comment:标题"`
	Message     string         `json:"message" gorm:"type:varchar(255);DEFAULT '';comment:消息"`
	Platform    string         `json:"platform" gorm:"type:varchar(20);DEFAULT '';comment:平台"`
	SendType    string         `json:"sendType" gorm:"type:varchar(20);DEFAULT '';comment:发送类型"`
	Delay       int            `json:"delay" gorm:"comment:延迟"`
	DelayNumber int            `json:"delayNumber" gorm:"comment:延迟数量"`
	DelayUnit   string         `json:"delayUnit" gorm:"type:varchar(10);DEFAULT '';comment:延迟单位"`
	Cron        string         `json:"cron" gorm:"type:varchar(15);DEFAULT '';comment:Cron"`
	LoopType    string         `json:"loopType" gorm:"type:varchar(20);DEFAULT '';comment:循环类型"`
	Day         string         `json:"day" gorm:"type:varchar(20);DEFAULT '';comment:日期"`
	Week        string         `json:"week" gorm:"type:varchar(30);DEFAULT '';comment:星期"`
	DayOfMonth  string         `json:"dayOfMonth" gorm:"type:varchar(20);DEFAULT '';comment:每月几号"`
	DayOfYear   string         `json:"dayOfYear" gorm:"type:varchar(20);DEFAULT '';comment:每年几月几号"`
	Time        string         `json:"time" gorm:"type:varchar(20);DEFAULT '';comment:时间"`
	Status      int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
}
