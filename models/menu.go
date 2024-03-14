package models

import (
	"errors"
	"gorm.io/gorm"
	"sort"
	"time"
)

// 菜单
type Menu struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Name       string         `json:"name" gorm:"type:varchar(255);DEFAULT '';comment:菜单/按钮名称"`
	Path       string         `json:"path" gorm:"type:varchar(255);DEFAULT '';comment:菜单路由地址"`
	Redirect   string         `json:"redirect" gorm:"type:varchar(255);DEFAULT '';comment:跳转路由地址"`
	Type       uint           `json:"type" gorm:"type:varchar(2);DEFAULT '';comment:菜单类型"`
	AuthCode   string         `json:"authcode" gorm:"type:varchar(255);DEFAULT '';comment:权限编码"`
	ParentID   uint           `json:"parentID" gorm:"type:varchar(64);DEFAULT '';comment:父节点ID"`
	Component  string         `json:"component" gorm:"type:varchar(40);DEFAULT '';comment:组件路径"`
	Sort       uint           `json:"sort" gorm:"comment:排序"`
	Title      string         `json:"title" gorm:"type:varchar(100);DEFAULT '';comment:标题"`
	Icon       string         `json:"icon" gorm:"type:varchar(20);DEFAULT '';comment:图标"`
	Hidden     bool           `json:"hidden" gorm:"type:tinyint(1);DEFAULT 0;comment:是否隐藏"`
	IsFrame    bool           `json:"isFrame" gorm:"type:tinyint(1);DEFAULT 0;comment:是否外链"`
	KeepAlive  bool           `json:"keepAlive" gorm:"type:tinyint(1);DEFAULT 0;comment:是否缓存"`
	FrameBlank bool           `json:"frameBlank" gorm:"type:tinyint(1);DEFAULT 0;comment:是否新窗口"`
	FrameSrc   string         `json:"frameSrc" gorm:"type:varchar(255);DEFAULT '';comment:外链地址"`
	Status     int            `json:"status" gorm:"type:tinyint(1);DEFAULT 0;comment:状态"`
	// Meta      string         `json:"meta" gorm:"type:varchar(255);DEFAULT '';comment:路由元数据"`
	// hidden isFrame keepAlive title
}

func AddMenu(name, path, autoCode, component, icon, frameSrc, title string, thetype, parentID, sort uint, hidden, isFrame, keepAlive, frameBlank bool) error {
	r := Menu{Name: name, Path: path, Type: thetype, Title: title, AuthCode: autoCode, ParentID: parentID, Component: component, Sort: sort, Icon: icon, Hidden: hidden, IsFrame: isFrame, KeepAlive: keepAlive, FrameBlank: frameBlank, FrameSrc: frameSrc}
	result := DB.Create(&r)
	err := result.Error
	return err
}

func UpdateMenu(id uint, updates map[string]interface{}) error {
	menu := Menu{}
	result := DB.Model(&menu).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteMenu(id uint) error {
	var menu Menu
	err := DB.Where("id = ?", id).Delete(&menu).Error
	return err
}

func GetMenuByID(menuID uint) (Menu, error) {
	var menu Menu
	result := DB.Where("id = ?", menuID).First(&menu)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return Menu{}, errors.New("menu_not_exist")
		} else {
			return Menu{}, result.Error
		}
	}
	return menu, nil
}

func GetAllMenus() (menus []Menu, err error) {
	err = DB.Find(&menus).Error
	return menus, err
}

type MenuMeta struct {
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Hidden     bool   `json:"hidden"`
	IsFrame    bool   `json:"isFrame"`
	KeepAlive  bool   `json:"keepAlive"`
	FrameBlank bool   `json:"frameBlank"`
	FrameSrc   string `json:"frameSrc"`
}

type MenuNode struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Redirect  string     `json:"redirect"`
	Type      uint       `json:"type"`
	AuthCode  string     `json:"authcode"`
	ParentID  uint       `json:"parentID"`
	Component string     `json:"component"`
	Sort      uint       `json:"sort"`
	Title     string     `json:"title"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuNode `json:"children"`
}

func buildMenuTree(menus map[uint]Menu, parentID uint) []MenuNode {
	var tree []MenuNode
	for _, menu := range menus {
		if menu.ParentID == parentID {
			meta := MenuMeta{
				Title:      menu.Title,
				Icon:       menu.Icon,
				Hidden:     menu.Hidden,
				IsFrame:    menu.IsFrame,
				KeepAlive:  menu.KeepAlive,
				FrameBlank: menu.FrameBlank,
				FrameSrc:   menu.FrameSrc,
			}
			node := MenuNode{
				Name:      menu.Name,
				Path:      menu.Path,
				Redirect:  menu.Redirect,
				Type:      menu.Type,
				AuthCode:  menu.AuthCode,
				ParentID:  menu.ParentID,
				Component: menu.Component,
				Sort:      menu.Sort,
				Title:     menu.Title,
				Meta:      meta,
				Children:  buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

func sortMenuTree(tree []MenuNode) {
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].Sort > tree[j].Sort
	})

	for _, node := range tree {
		sortMenuTree(node.Children)
	}
}

func GetAdminMenu(adminID uint) ([]MenuNode, error) {
	// 根据管理员 ID 获取管理员数据，假设这是从数据库中获取的管理员数据
	var admin Admin
	if err := DB.Preload("Roles.Menus").First(&admin, adminID).Error; err != nil {
		return nil, err
	}

	// 去重菜单
	uniqueMenus := make(map[uint]Menu)
	for _, role := range admin.Roles {
		for _, menu := range role.Menus {
			uniqueMenus[menu.ID] = menu
		}
	}

	// 创建一个映射来记录已存在的菜单项
	existingMenus := make(map[uint]bool)
	for _, menu := range uniqueMenus {
		existingMenus[menu.ID] = true
	}

	// 遍历菜单列表，检查父菜单是否存在于 uniqueMenus 中，如果不存在则添加到数组中
	for _, menu := range uniqueMenus {
		parentID := menu.ParentID
		// 检查父菜单是否存在于 uniqueMenus 中
		if parentID != 0 && !existingMenus[parentID] {
			parentMenu, err := GetMenuByID(parentID)
			if err != nil {
				return nil, err
			}
			// 将父菜单添加到 uniqueMenus 中
			uniqueMenus[parentID] = parentMenu
			existingMenus[parentID] = true
		}
	}

	// 构建菜单树状结构
	menuTree := buildMenuTree(uniqueMenus, 0)
	sortMenuTree(menuTree)
	return menuTree, nil
}
