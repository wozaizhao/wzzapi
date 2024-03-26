package models

import (
	"gorm.io/gorm"
	"time"
)

// 用户
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	NickName    string         `json:"nickName" gorm:"type:varchar(100);DEFAULT '';comment:昵称"`
	PhoneNumber string         `json:"phoneNumber" gorm:"unique type:varchar(20);DEFAULT '';comment:手机号"`
	Email       string         `json:"email" gorm:"type:varchar(40);DEFAULT '';comment:电子邮件"`
	Password    string         `json:"-" gorm:"type:varchar(100);DEFAULT '';comment:密码"`
	Avatar      string         `json:"avatar" gorm:"type:varchar(255);DEFAULT '';comment:头像网址"`
	Status      int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
	// OpenID      string         `json:"openID" gorm:"unique type:varchar(40);DEFAULT '';comment:openID"`

}
