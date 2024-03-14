package models

import (
	"errors"
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
	Roles       []Role         `json:"roles" gorm:"many2many:admin_roles;"`
	Permissions []string       `json:"permissions" gorm:"-"`
	Status      int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
}

func CreateAdmin(account, phoneNumber, email, password string) error {
	admin := Admin{Account: account, PhoneNumber: phoneNumber, Email: email, Password: password}
	result := DB.Create(&admin)
	err := result.Error
	return err
}

func VerifiyAdmin(account, phoneNumber, email, password string) (*Admin, error) {
	var admin Admin
	var err error
	if account != "" {
		err = DB.Where("account = ?", account).First(&admin).Error
	} else if phoneNumber != "" {
		err = DB.Where("phone_number = ?", phoneNumber).First(&admin).Error
	} else if email != "" {
		err = DB.Where("email = ?", email).First(&admin).Error
	} else {
		return nil, errors.New("must_provide_one_of_account_email_phone_number")
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("admin_not_exist")
		} else {
			return nil, err
		}
	}
	if admin.Password != password {
		return nil, errors.New("password_not_match")
	}
	return &admin, nil
}

func GetAdminByID(adminID uint) (*Admin, error) {
	var admin Admin
	result := DB.Where("id = ?", adminID).First(&admin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("admin_not_exist")
		} else {
			return nil, result.Error
		}
	}
	return &admin, nil
}

func UpdateAdmin(id uint, updates map[string]interface{}) error {
	admin := Admin{}
	result := DB.Model(&admin).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateAdminRoles(adminID uint, roleIDs []uint) error {
	admin := Admin{ID: adminID}
	err := DB.First(&admin).Error
	if err != nil {
		return err
	}

	var roles []Role
	err = DB.Find(&roles, roleIDs).Error
	if err != nil {
		return err
	}

	err = DB.Model(&admin).Association("Roles").Replace(&roles)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAdmin(id uint) error {
	var admin Admin
	err := DB.Where("id = ?", id).Delete(&admin).Error
	return err
}

func GetAdmins(pageNum, pageSize int, keyword string) (admins []Admin, err error) {
	db := DB
	if keyword != "" {
		db = db.Scopes(Search(keyword, "account", "phone_number", "email"))
	}
	if pageNum != 0 && pageSize != 0 {
		db = db.Scopes(Paginate(pageNum, pageSize))
	}
	err = db.Order("created_at asc").Preload("Roles").Find(&admins).Error
	return admins, err
}

func GetAdminsCount(keyword string) int64 {
	var count int64
	db := DB
	if keyword != "" {
		db = db.Scopes(Search(keyword, "account", "phone_number", "email"))
	}
	db.Model(&Admin{}).Count(&count)
	return count
}
