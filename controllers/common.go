package controllers

import (
	"github.com/gin-gonic/gin"
	"wozaizhao.com/wzzapi/services"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		RenderBadRequest(c, err)
		return
	}
	dir := c.DefaultPostForm("dir", "wzz")
	uploadedURL, err := services.UploadByFile(dir, file)

	if err != nil {
		RenderFail(c, err.Error())
		return
	}

	RenderSuccess(c, uploadedURL, "upload_success")

}
