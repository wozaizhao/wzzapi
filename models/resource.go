package models

import (
	"time"

	"gorm.io/gorm"
)

// 资源 一个个的网址，其实是一个程序或一些分享的资源
type Resource struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	Title          string         `json:"title" gorm:"type:varchar(100);DEFAULT '';comment:名称"`
	Url            string         `json:"url" gorm:"type:varchar(255);DEFAULT '';comment:网址"`
	Logo           string         `json:"logo" gorm:"type:varchar(255);DEFAULT '';comment:Logo"`
	Comment        string         `json:"comment" gorm:"type:varchar(255);DEFAULT '';comment:评论"`
	Visible        bool           `json:"visible" gorm:"type:tinyint(1);DEFAULT 0;comment:是否可见"`
	Tags           string         `json:"tags" gorm:"type:varchar(255);DEFAULT '';comment:标签"`
	Sort           uint           `json:"sort" gorm:"comment:排序"`
	ClickCount     int            `json:"clickCount" gorm:"type:int;DEFAULT 0;comment:点击数"`
	CreatedByAdmin uint           `json:"createdByAdmin"`
	CreatedByUser  uint           `json:"createdByUser"`
}

func (Resource) TableName() string {
	return "business_resource"
}

func CreateResource(title, logo, url, comment, tags string, creator uint) error {
	resource := Resource{Title: title, Logo: logo, Url: url, Comment: comment, Tags: tags, CreatedByAdmin: creator}
	result := DB.Create(&resource)
	err := result.Error
	return err
}

func UpdateResource(id uint, updates map[string]interface{}) error {
	resource := Resource{}
	result := DB.Model(&resource).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 设置资源可见
func SetResourceVisible(id uint, visible bool) error {
	resource := Resource{}
	result := DB.Model(&resource).Where("id = ?", id).Update("visible", visible)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetResources(pageNum, pageSize int, tag string) (resources []Resource, err error) {
	db := DB
	if tag != "" {
		db = db.Scopes(Search(tag, "tags"))
	}
	if pageNum != 0 && pageSize != 0 {
		db = db.Scopes(Paginate(pageNum, pageSize))
	}
	err = db.Order("sort desc,click_count desc").Find(&resources).Error
	return resources, err
}

func GetResourcesCount(tag string) int64 {
	var count int64
	db := DB
	if tag != "" {
		db = db.Scopes(Search(tag, "tags"))
	}
	db.Model(&Resource{}).Count(&count)
	return count
}

func DeleteResource(id uint) error {
	var resource Resource
	err := DB.Where("id = ?", id).Delete(&resource).Error
	return err
}
