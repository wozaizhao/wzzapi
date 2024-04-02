package models

import (
	"time"

	"gorm.io/gorm"
)

// 字典
type Dict struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Label     string         `json:"label" gorm:"type:varchar(50);DEFAULT '';comment:标识"`
	Value     string         `json:"value" gorm:"type:varchar(50);DEFAULT '';comment:值"`
	DictName  string         `json:"dictName" gorm:"type:varchar(50);DEFAULT '';comment:字典名称"`
	DictType  string         `json:"dictType" gorm:"type:varchar(50);DEFAULT '';comment:字典类型"`
	Remark    string         `json:"remark" gorm:"type:varchar(100);DEFAULT '';comment:备注"`
	Status    int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
}

func (Dict) TableName() string {
	return "sys_dicts"
}

func CreateDict(label, value, dictName, dictType, remark string) (dict Dict, err error) {
	dict = Dict{Label: label, Value: value, DictName: dictName, DictType: dictType, Remark: remark}
	result := DB.Create(&dict)
	err = result.Error
	return dict, err
}

func GetDictsByType(dictType string) (dicts []Dict, err error) {
	err = DB.Where("dict_type = ?", dictType).Find(&dicts).Error
	return dicts, err
}

func UpdateDict(id uint, label, value, remark string) error {
	dict := Dict{}
	updates := map[string]interface{}{
		"label":  label,
		"value":  value,
		"remark": remark,
	}
	result := DB.Model(&dict).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteDict(id uint) error {
	var dict Dict
	err := DB.Where("id = ?", id).Delete(&dict).Error
	return err
}

func GetDicts(thetype string, keyword string) (dicts []Dict, err error) {
	if keyword != "" {
		err = DB.Scopes(Search(keyword, "dict_name", "dict_type")).Order("created_at asc").Find(&dicts).Error
	} else {
		err = DB.Order("created_at asc").Find(&dicts).Error
	}
	return dicts, err
}
