package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"wozaizhao.com/wzzapi/global"
	"wozaizhao.com/wzzapi/notify/wecom"
)

// 通知
type NotifyWecom struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	wecom.NotifyConfig
	CreatedBy uint   `json:"createdBy" gorm:"type:varchar(64);DEFAULT '';comment:创建人"`
	Remark    string `json:"remark" gorm:"type:varchar(255);DEFAULT '';comment:备注"`
}

func (NotifyWecom) TableName() string {
	return "notify_wecom"
}

func (n *NotifyWecom) BeforeSave(tx *gorm.DB) (err error) {
	n.WebhookURL = encrypt(n.WebhookURL)
	return nil
}

func (n *NotifyWecom) AfterFind(tx *gorm.DB) error {
	webhookURL := decrypt(n.WebhookURL)
	n.WebhookURL = global.MaskSensitiveInfo(webhookURL, len(webhookURL)-20, 20, "*")
	return nil
}

// 增加
func CreateNotifyWecom(name, webhookURL, remark string, userID uint, dry bool) error {
	notify := NotifyWecom{Remark: remark, CreatedBy: userID}
	notify.WebhookURL = webhookURL
	notify.NotifyName = name
	notify.Dry = dry
	result := DB.Create(&notify)
	err := result.Error
	return err
}

// 删除
func DeleteNotifyWecom(id, userID uint) error {
	var notify NotifyWecom
	err := DB.Where("id = ? AND created_by = ?", id, userID).Delete(&notify).Error
	return err
}

// 更新
func UpdateNotifyWecom(id uint, updates map[string]interface{}) error {
	notify := NotifyWecom{}
	result := DB.Model(&notify).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 查询列表
func GetNotifyWecomList(userID uint) (nitifies []NotifyWecom, err error) {
	err = DB.Where("created_by = ?", userID).Find(&nitifies).Error
	return
}

// 查询单个
func GetNotifyWecomByID(id, userID uint) (*NotifyWecom, error) {
	var notify NotifyWecom
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
