package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"wozaizhao.com/wzzapi/global"
	"wozaizhao.com/wzzapi/models"
	"wozaizhao.com/wzzapi/notify/dingtalk"
	"wozaizhao.com/wzzapi/notify/wecom"
)

type platform struct {
	Platform string `json:"platform" binding:"required"`
}

type adminAddNotifyParam struct {
	Name       string `json:"name" binding:"required"`
	WebHookURL string `json:"webHookURL" binding:"required"`
	SignSecret string `json:"signSecret"`
	Remark     string `json:"remark"`
	Dry        bool   `json:"dry" binding:"required"`
	platform
}

func AdminAddNotify(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	var req adminAddNotifyParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	if req.Platform == "dingtalk" {
		err := models.CreateNotifyDingtalk(req.Name, req.WebHookURL, req.SignSecret, req.Remark, adminID, req.Dry)
		if err != nil {
			RenderError(c, err)
			return
		}
	} else if req.Platform == "wecom" {
		err := models.CreateNotifyWecom(req.Name, req.WebHookURL, req.Remark, adminID, req.Dry)
		if err != nil {
			RenderError(c, err)
			return
		}
	}

	RenderSuccess(c, "", "add_notify_success")
}

func AdminDeleteNotify(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	var req idParam
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	var param platform
	if err := c.BindJSON(&param); err != nil {
		RenderBadRequest(c, err)
		return
	}
	if param.Platform == "dingtalk" {
		err := models.DeleteNotifyDingtalk(req.ID, adminID)
		if err != nil {
			RenderError(c, err)
			return
		}
	} else if param.Platform == "wecom" {
		err := models.DeleteNotifyWecom(req.ID, adminID)
		if err != nil {
			RenderError(c, err)
			return
		}
	}
	RenderSuccess(c, "", "delete_notify_success")
}

func AdminUpdateNotify(c *gin.Context) {

}

func AdminGetNotifies(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)

	dingtalks, err := models.GetNotifyDingtalkList(adminID)
	if err != nil {
		RenderError(c, err)
		return
	}
	wecoms, err := models.GetNotifyWecomList(adminID)
	if err != nil {
		RenderError(c, err)
		return
	}
	RenderSuccess(c, gin.H{"dingtalks": dingtalks, "wecoms": wecoms}, "get_notifies_success")
}

func AdminGetNotify(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	var req idParam
	if err := c.ShouldBindUri(&req); err != nil {
		RenderError(c, err)
		return
	}
	var param platform
	if err := c.BindJSON(&param); err != nil {
		RenderBadRequest(c, err)
		return
	}

	if param.Platform == "dingtalk" {
		notify, err := models.GetNotifyDingtalkByID(req.ID, adminID)
		if err != nil {
			RenderError(c, err)
			return
		}
		RenderSuccess(c, notify, "get_notify_success")
	} else if param.Platform == "wecom" {
		notify, err := models.GetNotifyWecomByID(req.ID, adminID)
		if err != nil {
			RenderError(c, err)
			return
		}
		RenderSuccess(c, notify, "get_notify_success")
	}
}

type adminSendMessageParam struct {
	Title   string `json:"title"`
	Message string `json:"message" binding:"required"`
	Dry     bool   `json:"dry"`
	Sender  uint   `json:"sender" binding:"required"`
	platform
}

func AdminSendMessage(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint)
	var req adminSendMessageParam
	if err := c.BindJSON(&req); err != nil {
		RenderBadRequest(c, err)
		return
	}
	if req.Platform == "dingtalk" {

		sender, err := models.GetNotifyDingtalkByID(req.Sender, adminID)
		if err != nil {
			RenderError(c, err)
			return
		}
		conf := &dingtalk.NotifyConfig{
			SignSecret: sender.SignSecret,
			WebhookURL: sender.WebhookURL,
		}
		conf.Dry = req.Dry
		err = conf.Config(global.NotifySettings{})
		if err != nil {
			RenderError(c, err)
			return
		}
		err = conf.SendDingtalkNotification(req.Title, req.Message)
		if err != nil {
			RenderError(c, err)
			return
		}
		RenderSuccess(c, "", "send_dingtalk_notification_success")
	} else if req.Platform == "wecom" {
		sender, err := models.GetNotifyWecomByID(req.Sender, adminID)
		if err != nil {
			RenderError(c, err)
			return
		}
		conf := &wecom.NotifyConfig{
			WebhookURL: sender.WebhookURL,
		}
		conf.Dry = req.Dry
		err = conf.Config(global.NotifySettings{})
		if err != nil {
			RenderError(c, err)
			return
		}
		err = conf.SendWecomNotification(req.Message)
		if err != nil {
			RenderError(c, err)
			return
		}
		RenderSuccess(c, "", "send_wecom_notification_success")
	} else {
		RenderBadRequest(c, errors.New("platform_is_empty"))
	}

}
