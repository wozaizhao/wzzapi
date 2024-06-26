package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"wozaizhao.com/wzzapi/global"
	"wozaizhao.com/wzzapi/notify/dingtalk"
)

// 通知
type NotifyDingtalk struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	dingtalk.NotifyConfig
	CreatedBy           uint   `json:"createdBy" gorm:"type:varchar(64);DEFAULT '';comment:创建人"`
	Remark              string `json:"remark" gorm:"type:varchar(255);DEFAULT '';comment:备注"`
	WebhookURLSensitive string `json:"webhook" gorm:"-"`
}

func (NotifyDingtalk) TableName() string {
	return "notify_dingtalk"
}

func (n *NotifyDingtalk) BeforeSave(tx *gorm.DB) (err error) {
	if n.SignSecret != "" {
		n.SignSecret = encrypt(n.SignSecret)
	}
	n.WebhookURL = encrypt(n.WebhookURL)
	return nil
}

func (n *NotifyDingtalk) AfterFind(tx *gorm.DB) error {
	if n.SignSecret != "" {
		cipherSign := decrypt(n.SignSecret)
		n.SignSecret = global.MaskSensitiveInfo(cipherSign, 5, 6, "*")
	}

	webhookURL := decrypt(n.WebhookURL)
	n.WebhookURL = webhookURL
	n.WebhookURLSensitive = global.MaskSensitiveInfo(webhookURL, len(webhookURL)-20, 20, "*")
	return nil
}

// 增加
func CreateNotifyDingtalk(name, webhookURL, signSecret, remark string, userID uint, dry bool) error {
	notify := NotifyDingtalk{Remark: remark, CreatedBy: userID}
	notify.WebhookURL = webhookURL
	notify.SignSecret = signSecret
	notify.NotifyName = name
	notify.Dry = dry
	result := DB.Create(&notify)
	err := result.Error
	return err
}

// 删除
func DeleteNotifyDingtalk(id, userID uint) error {
	var notify NotifyDingtalk
	err := DB.Where("id = ? AND created_by = ?", id, userID).Delete(&notify).Error
	return err
}

// 更新
func UpdateNotifyDingtalk(id uint, updates map[string]interface{}) error {
	notify := NotifyDingtalk{}
	result := DB.Model(&notify).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 查询列表
func GetNotifyDingtalkList(userID uint) (nitifies []NotifyDingtalk, err error) {
	err = DB.Where("created_by = ?", userID).Find(&nitifies).Error
	return
}

// 查询单个
func GetNotifyDingtalkByID(id, userID uint) (*NotifyDingtalk, error) {
	var notify NotifyDingtalk
	result := DB.Where("id = ? and created_by = ?", id, userID).First(&notify)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("notify_not_exist")
		} else {
			return nil, result.Error
		}
	}
	return &notify, nil
}
