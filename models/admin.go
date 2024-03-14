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
	Permissions []string       `json:"permissions" gorm:"-"`
	Status      int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
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
