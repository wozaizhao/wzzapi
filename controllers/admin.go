package controllers

import (
	"github.com/gin-gonic/gin"
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
