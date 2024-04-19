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
	r.GET("/dicts", controllers.GetDictsByType)

	// 注册用户可以访问 /user
	user := r.Group("/user", controllers.UserAuth())
	{
		user.GET("/current")
	}

	admin := r.Group("/admin", controllers.AdminAuth())
	{
		admin.GET("/current", controllers.GetCurrentAdmin)
		admin.POST("/upload", controllers.Upload)
		// 增加、删除、修改、查询PC端用户
		admin.POST("/admins", controllers.AdminAddAdmin)
		admin.DELETE("/admins/:id", controllers.AdminDeleteAdmin)
		admin.PUT("/admins", controllers.AdminUpdateAdmin)
		admin.GET("/admins", controllers.AdminGetAdmins)
		admin.POST("/admins/toggleStatus", controllers.ToggleAdminStatus)
		// 增加、删除、修改、查询菜单
		admin.GET("/menuList", controllers.GetMenuList)
		admin.POST("/menus", controllers.AdminAddMenu)
		admin.DELETE("/menus/:id", controllers.AdminDeleteMenu)
		admin.PUT("/menus", controllers.AdminUpdateMenu)
		admin.GET("/menus", controllers.AdminGetMenus)
		admin.GET("/menus/:id", controllers.AdminGetMenu)
		// 增加、删除、修改、查询角色
		admin.POST("/roles", controllers.AdminAddRole)
		admin.DELETE("/roles/:id", controllers.AdminDeleteRole)
		admin.PUT("/roles", controllers.AdminUpdateRole)
		admin.GET("/roles", controllers.AdminGetRoles)
		admin.GET("/roles/:id/menus", controllers.AdminGetRoleMenus)
		admin.GET("/allRoles", controllers.AdminGetAllRoles)
		// 增加、删除、修改、查询notify
		admin.POST("/notifies", controllers.AdminAddNotify)
		admin.DELETE("/notifies/:id", controllers.AdminDeleteNotify)
		admin.PUT("/notifies", controllers.AdminUpdateNotify)
		admin.GET("/notifies", controllers.AdminGetNotifies)
		admin.GET("/notifies/:id", controllers.AdminGetNotify)
		// 即时发送消息
		admin.POST("/send", controllers.AdminSendMessage)
	}

	return r
}
