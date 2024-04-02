package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"wozaizhao.com/wzzapi/global"
	"wozaizhao.com/wzzapi/models"
)

type adminAddRoleParam struct {
	Name    string `json:"name" binding:"required"`
	Remark  string `json:"remark"`
	MenuIDs []uint `json:"menus"`
}

func AdminAddRole(c *gin.Context) {
	var req adminAddRoleParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	err := models.AddRole(req.Name, req.Remark, req.MenuIDs)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "add_role_success")
}

type adminUpdateRoleParam struct {
	ID      uint   `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Remark  string `json:"remark"`
	MenuIDs []uint `json:"menus"`
}

func AdminUpdateRole(c *gin.Context) {
	var req adminUpdateRoleParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	err := models.UpdateRole(req.ID, req.Name, req.Remark, req.MenuIDs)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "update_role_success")
}

type idParam struct {
	ID uint `uri:"id" binding:"required"`
}

func AdminDeleteRole(c *gin.Context) {
	var req idParam
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	err := models.DeleteRole(req.ID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "delete_role_success")
}

func AdminGetRoles(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	pageNumParam := c.DefaultQuery("pageNum", "1")
	pageSizeParam := c.DefaultQuery("pageSize", "10")
	pageNum, _ := global.ParseInt(pageNumParam)
	pageSize, _ := global.ParseInt(pageSizeParam)
	roles, err := models.GetRoles(int(pageNum), int(pageSize), keyword)

	if err != nil {
		RenderError(c, err)
		return
	}
	count := models.GetRolesCount(keyword)
	var res commonList
	res.List = roles
	res.Total = count
	RenderSuccess(c, res, "get_roles_success")
}

type adminAddMenuParam struct {
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Type       uint   `json:"type"`
	AuthCode   string `json:"authCode"`
	ParentID   uint   `json:"parentID"`
	Component  string `json:"component" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Sort       uint   `json:"sort"`
	Icon       string `json:"icon"`
	Hidden     bool   `json:"hidden"`
	IsFrame    bool   `json:"isFrame"`
	KeepAlive  bool   `json:"keepAlive"`
	FrameBlank bool   `json:"frameBlank"`
	FrameSrc   string `json:"frameSrc"`
	// Meta      string `json:"meta"`
}

func AdminAddMenu(c *gin.Context) {
	var req adminAddMenuParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	err := models.AddMenu(req.Name, req.Path, req.AuthCode, req.Component, req.Icon, req.FrameSrc, req.Title, req.Type, req.ParentID, req.Sort, req.Hidden, req.IsFrame, req.KeepAlive, req.FrameBlank)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "add_menu_success")
}

type adminUpdateMenuParam struct {
	ID         uint   `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Type       uint   `json:"type"`
	AuthCode   string `json:"authCode"`
	ParentID   uint   `json:"parentID"`
	Component  string `json:"component" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Sort       uint   `json:"sort"`
	Icon       string `json:"icon"`
	Hidden     bool   `json:"hidden"`
	IsFrame    bool   `json:"isFrame"`
	KeepAlive  bool   `json:"keepAlive"`
	FrameBlank bool   `json:"frameBlank"`
	FrameSrc   string `json:"frameSrc"`
	// Meta      string `json:"meta"`
}

func AdminUpdateMenu(c *gin.Context) {
	var req adminUpdateMenuParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"path":        req.Path,
		"type":        req.Type,
		"auth_code":   req.AuthCode,
		"parent_id":   req.ParentID,
		"component":   req.Component,
		"title":       req.Title,
		"sort":        req.Sort,
		"icon":        req.Icon,
		"frame_src":   req.FrameSrc,
		"hidden":      req.Hidden,
		"is_frame":    req.IsFrame,
		"keep_alive":  req.KeepAlive,
		"frame_blank": req.FrameBlank,
		// "meta":      req.Meta,
	}
	err := models.UpdateMenu(req.ID, updates)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "update_menu_success")
}

type TreeNode struct {
	Label        string     `json:"label"`
	Value        uint       `json:"value"`
	ParentName   *string    `json:"parentName,omitempty"`
	AssociatedID *uint      `json:"associatedID,omitempty"`
	Children     []TreeNode `json:"children"`
}

type OriginalNode struct {
	Label        string
	Value        uint
	ParentID     uint
	AssociatedID *uint
}

func getTree(nodes []OriginalNode, parentID uint, parentName string) []TreeNode {
	var tree []TreeNode
	for _, node := range nodes {
		if node.ParentID == parentID {
			node := TreeNode{
				Label:        node.Label,
				Value:        node.Value,
				ParentName:   &parentName,
				AssociatedID: node.AssociatedID,
				Children:     getTree(nodes, node.Value, node.Label),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// 管理后台获取菜单（树型结构）
func AdminGetMenus(c *gin.Context) {
	menus, err := models.GetAllMenus()
	if err != nil {
		RenderError(c, err)
		return
	}
	var pureMenus []OriginalNode
	for _, menu := range menus {
		item := OriginalNode{
			Label:    menu.Title,
			Value:    menu.ID,
			ParentID: menu.ParentID,
		}
		pureMenus = append(pureMenus, item)
	}
	// 获取根节点菜单
	var rootMenus []OriginalNode
	for _, menu := range pureMenus {
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
		}
	}

	// 构建树状结构
	var menuTree []TreeNode
	for _, rootMenu := range rootMenus {
		node := TreeNode{
			Label:    rootMenu.Label,
			Value:    rootMenu.Value,
			Children: getTree(pureMenus, rootMenu.Value, ""),
		}
		menuTree = append(menuTree, node)
	}
	RenderSuccess(c, menuTree, "get_menu_tree_success")
}

// 管理后台获取单个菜单
func AdminGetMenu(c *gin.Context) {
	var req idParam
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	menu, err := models.GetMenuByID(req.ID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, menu, "get_menu_success")
}

// 管理后台删除单个菜单
func AdminDeleteMenu(c *gin.Context) {
	var req idParam
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	err := models.DeleteMenu(req.ID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "delete_menu_success")
}

func AdminGetRoleMenus(c *gin.Context) {
	var req idParam
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	menus, err := models.GetRoleMenus(req.ID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, menus, "get_role_menus_success")

}

// 管理员获取自己的菜单（根据角色权限）
func GetMenuList(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	menuList, err := models.GetAdminMenu(adminID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, menuList, "get_admin_menu_list_success")
}

func AdminGetAllRoles(c *gin.Context) {
	roles, err := models.GetAllRoles()
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, roles, "get_all_roles_success")
}

func GetDictsByType(c *gin.Context) {
	dictType := c.DefaultQuery("dictType", "")
	if dictType != "" {
		dicts, err := models.GetDictsByType(dictType)
		if err != nil {
			RenderError(c, err)
			return
		}
		RenderSuccess(c, dicts, "get_dicts_by_type_success")
	} else {
		RenderError(c, errors.New("no_type_found"))
	}
}
