package controllers

import (
	"github.com/gin-gonic/gin"
	"wozaizhao.com/wzzapi/global"
	"wozaizhao.com/wzzapi/models"
)

type adminLoginByPasswordReq struct {
	Account     string `json:"account"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
}

func AdminLoginByPassword(c *gin.Context) {
	var req adminLoginByPasswordReq
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	admin, err := models.VerifiyAdmin(req.Account, req.PhoneNumber, req.Email, req.Password)
	if err != nil || admin == nil {
		RenderError(c, err)
		return
	}

	token, err := generateToken(0, admin.ID)
	if err != nil {
		RenderError(c, err)
		return
	}

	RenderSuccess(c, token, "login_success")
}

func GetCurrentAdmin(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	admin, err := models.GetAdminByID(adminID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, admin, "get_admininfo_success")
}

type adminAddAdminParam struct {
	Account     string `json:"account"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
}

func AdminAddAdmin(c *gin.Context) {
	var req adminAddAdminParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}

	err := models.CreateAdmin(req.Account, req.PhoneNumber, req.Email, req.Password)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "add_admin_success")
}

type adminUpdateAdminParam struct {
	ID          uint   `json:"id" binding:"required"`
	Email       string `json:"email" `
	PhoneNumber string `json:"phoneNumber"`
	Roles       []uint `json:"roles"`
}

func AdminUpdateAdmin(c *gin.Context) {
	var req adminUpdateAdminParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	updates := map[string]interface{}{
		"email":        req.Email,
		"phone_number": req.PhoneNumber,
	}
	err := models.UpdateAdmin(req.ID, updates)
	errUpdateRole := models.UpdateAdminRoles(req.ID, req.Roles)
	if err != nil {
		RenderError(c, err)
		return
	}
	if errUpdateRole != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "update_admin_success")
}

// type adminDeleteAdminParam struct {
// 	ID uint `uri:"id" binding:"required"`
// }

func AdminDeleteAdmin(c *gin.Context) {
	var req IDInUri
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	err := models.DeleteAdmin(req.ID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "delete_admin_success")
}

func AdminGetAdmins(c *gin.Context) {
	// adminID := c.MustGet("adminID").(uint)
	keyword := c.DefaultQuery("keyword", "")
	pageNumParam := c.DefaultQuery("pageNum", "1")
	pageSizeParam := c.DefaultQuery("pageSize", "10")
	pageNum, _ := global.ParseInt(pageNumParam)
	pageSize, _ := global.ParseInt(pageSizeParam)
	admins, err := models.GetAdmins(int(pageNum), int(pageSize), keyword)

	if err != nil {
		RenderError(c, err)
		return
	}
	count := models.GetAdminsCount(keyword)
	var res commonList
	res.List = admins
	res.Total = count
	RenderSuccess(c, res, "get_admins_success")
}

type toggleStatusParam struct {
	ID     uint `json:"id" binding:"required"`
	Status uint `json:"status" `
}

func ToggleAdminStatus(c *gin.Context) {
	var req toggleStatusParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	updates := map[string]interface{}{
		"status": req.Status,
	}
	err := models.UpdateAdmin(req.ID, updates)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "toggle_admin_status_success")
}
