package models

import (
	"gorm.io/gorm"
	"time"
)

// 管理员,后台用户
type Admin struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Account     string         `json:"account" gorm:"unique type:varchar(8);DEFAULT '';comment:帐号名"`
	Password    string         `json:"-" gorm:"type:varchar(255);DEFAULT '';comment:密码"`
	Avatar      string         `json:"avatar" gorm:"type:varchar(255);DEFAULT '';comment:头像网址"`
	Email       string         `json:"email" gorm:"unique type:varchar(255);DEFAULT '';comment:email"`
	PhoneNumber string         `json:"phoneNumber" gorm:"unique type:varchar(20);DEFAULT '';comment:手机号"`
	Status      int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
}
