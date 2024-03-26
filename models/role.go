package models

import (
	log "github.com/sirupsen/logrus"
	"time"

	"gorm.io/gorm"
)

// 角色
type Role struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"name" gorm:"unique type:varchar(40);DEFAULT '';comment:角色名称"`
	Remark    string         `json:"remark" gorm:"type:varchar(100);DEFAULT '';comment:备注"`
	Menus     []Menu         `json:"menus" gorm:"many2many:sys_role_menus;"`
	Status    int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
}

func (Role) TableName() string {
	return "sys_roles"
}

func AddRole(name, remark string, menuIDs []uint) error {
	r := Role{Name: name, Remark: remark}
	menus, err := GetAllMenus()
	if err != nil {
		return err
	}
	// 将菜单与角色关联
	for _, menuID := range menuIDs {
		for _, menu := range menus {
			if menu.ID == menuID {
				r.Menus = append(r.Menus, menu)
				break
			}
		}
	}
	result := DB.Create(&r)
	err = result.Error
	return err
}

func UpdateRole(id uint, name, remark string, menuIDs []uint) error {
	// 根据 ID 获取角色，假设这是从数据库中获取的角色数据
	var role Role
	if err := DB.Preload("Menus").First(&role, id).Error; err != nil {
		return err
	}

	// 清除旧的关联数据
	errClear := DB.Model(&role).Association("Menus").Delete(role.Menus)
	if errClear != nil {
		log.Errorf("UpdateRole Failed: %s", errClear)
		return errClear
	}

	// 更新角色的属性
	role.Name = name
	role.Remark = remark

	// 清空菜单列表
	role.Menus = []Menu{}

	menus, err := GetAllMenus()
	if err != nil {
		return err
	}

	// 将菜单与角色关联
	for _, menuID := range menuIDs {
		for _, menu := range menus {
			if menu.ID == menuID {
				role.Menus = append(role.Menus, menu)
				break
			}
		}
	}

	// 保存更新后的角色到数据库
	if err := DB.Save(&role).Error; err != nil {
		return err
	}

	return nil
}

func DeleteRole(id uint) error {
	var role Role
	err := DB.Where("id = ?", id).Delete(&role).Error
	return err
}

func GetRoles(pageNum, pageSize int, keyword string) (roles []Role, err error) {
	if keyword != "" {
		err = DB.Scopes(Paginate(pageNum, pageSize), Search(keyword, "name", "remark")).Order("created_at asc").Find(&roles).Error
	} else {
		err = DB.Scopes(Paginate(pageNum, pageSize)).Order("created_at asc").Find(&roles).Error
	}
	return roles, err
}

func GetRolesCount(keyword string) int64 {
	var count int64
	DB.Scopes(Search(keyword, "name", "remark")).Model(&Role{}).Count(&count)
	return count
}

func GetAllRoles() (roles []Role, err error) {
	err = DB.Find(&roles).Error
	return roles, err
}

func GetRoleMenus(roleID uint) ([]Menu, error) {
	// 根据角色 ID 获取角色，假设这是从数据库中获取的角色数据
	var role Role
	if err := DB.Preload("Menus").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	return role.Menus, nil
}
