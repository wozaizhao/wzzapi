package router

import (
	"github.com/gin-gonic/gin"
	"wozaizhao.com/wzzapi/controllers"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	// r.SetTrustedProxies([]string{"192.168.1.2"})

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	// allows all origins
	// r.Use(cors.Default())

	// 任何人都可访问
	r.POST("/login")
	r.POST("/adminLogin", controllers.AdminLoginByPassword)

	// 注册用户可以访问 /user
	user := r.Group("/user", controllers.UserAuth())
	{
		user.GET("/current")
	}

	admin := r.Group("/admin", controllers.AdminAuth())
	{
		admin.GET("/current", controllers.GetCurrentAdmin)
	}

	return r
}
