package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"wozaizhao.com/wzzapi/models"
)

// get resources
func GetResources(c *gin.Context) {
	// userID := c.MustGet("userID").(uint)
	pageParam := c.DefaultQuery("page", "1")
	pageSizeParam := c.DefaultQuery("pageSize", "10")
	tagParam := c.DefaultQuery("tag", "")

	page, _ := strconv.Atoi(pageParam)
	pageSize, _ := strconv.Atoi(pageSizeParam)
	resources, err := models.GetResources(page, pageSize, tagParam)
	if err != nil {
		RenderFail(c, err.Error())
		return
	}
	count := models.GetResourcesCount(tagParam)
	var res commonList
	res.List = resources
	res.Total = count
	RenderSuccess(c, res, "get_resources_success")
}

type Resource struct {
	ID      uint   `json:"id"`
	Title   string `json:"title" binding:"required"`
	Url     string `json:"url" binding:"required"`
	Logo    string `json:"logo"`
	Tags    string `json:"tags"`
	Comment string `json:"comment"`
}

func AdminUpdateResource(c *gin.Context) {
	var resource Resource
	if err := c.ShouldBindJSON(&resource); err != nil {
		RenderBadRequest(c, err)
		return
	}
	if resource.ID == 0 {
		RenderBadRequest(c, errors.New("id_is_required"))
		return
	}
	updates := map[string]interface{}{
		"Title":   resource.Title,
		"Url":     resource.Url,
		"Logo":    resource.Logo,
		"Comment": resource.Comment,
	}
	err := models.UpdateResource(resource.ID, updates)
	if err != nil {
		RenderFail(c, err.Error())
		return
	}
	RenderSuccess(c, nil, "update_success")
}

type setVisibleReqParam struct {
	ID      uint `json:"id"`
	Visible bool `json:"visible"`
}

func AdminSetResourceVisible(c *gin.Context) {
	var resource setVisibleReqParam
	if err := c.ShouldBindJSON(&resource); err != nil {
		RenderBadRequest(c, err)
		return
	}
	err := models.SetResourceVisible(resource.ID, resource.Visible)
	if err != nil {
		RenderFail(c, err.Error())
		return
	}
	RenderSuccess(c, nil, "update_success")
}

func AdminAddResource(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	var r Resource
	if err := c.ShouldBindJSON(&r); err != nil {
		RenderBadRequest(c, err)
		return
	}
	err := models.CreateResource(r.Title, r.Logo, r.Url, r.Comment, r.Tags, adminID)
	if err != nil {
		RenderFail(c, err.Error())
		return
	}
	RenderSuccess(c, nil, "add_success")
}

func AdminDeleteResource(c *gin.Context) {
	var req IDInUri
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	err := models.DeleteResource(req.ID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, "", "delete_resource_success")
}
