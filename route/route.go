package route

import (
	"embed"
	"io"
	"io/fs"
	"msm/api"
	"msm/config"
	"msm/consts/permission"
	"msm/consts/role"
	"msm/log"
	"msm/route/middle"
	"msm/utils"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.Default()
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.DebugMode)
	routePathInit(r)
	staticInit(r)
	pprofInit(r)
	r.Run(config.CF.Listen)
}

//go:embed templates
var f embed.FS

func staticInit(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		b, _ := f.ReadFile("templates/index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})
	r.StaticFS("/js", http.FS(utils.UnwarpIgnore(fs.Sub(f, "templates/js"))))
	r.StaticFS("/css", http.FS(utils.UnwarpIgnore(fs.Sub(f, "templates/css"))))
	r.StaticFS("/media", http.FS(utils.UnwarpIgnore(fs.Sub(f, "templates/media"))))
	r.StaticFS("/fonts", http.FS(utils.UnwarpIgnore(fs.Sub(f, "templates/fonts"))))
}

func pprofInit(r *gin.Engine) {
	if config.CF.PprofEnable {
		pprof.Register(r)
		log.Logger.Info("启用 pprof")
	}
}

func routePathInit(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use(middle.CheckToken())
	apiGroup.Use(middle.PanicMiddle())
	{
		apiGroup.GET("/ws", middle.OprPermission(permission.TERMINAL_OPERATION), api.WsApi.WebsocketHandle)

		processGroup := apiGroup.Group("/process")
		{
			processGroup.DELETE("", middle.OprPermission(permission.STOP_OPERATION), api.ProcApi.KillProcess)
			processGroup.GET("", api.ProcApi.GetProcessList)
			processGroup.PUT("", middle.OprPermission(permission.START_OPERATION), api.ProcApi.StartProcess)
			processGroup.GET("/control", middle.RolePermission(role.ADMIN), api.ProcApi.ProcessControl)

			proConfigGroup := processGroup.Group("/config")
			{
				proConfigGroup.POST("", middle.RolePermission(role.ROOT), api.ProcApi.CreateNewProcess)
				proConfigGroup.DELETE("", middle.RolePermission(role.ROOT), api.ProcApi.DeleteNewProcess)
				proConfigGroup.PUT("", middle.RolePermission(role.ROOT), api.ProcApi.UpdateProcessConfig)
				proConfigGroup.GET("", middle.RolePermission(role.ADMIN), api.ProcApi.GetProcessConfig)
			}
		}

		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/login", api.UserApi.LoginHandler)
			userGroup.POST("", middle.RolePermission(role.ROOT), api.UserApi.CreateUser)
			userGroup.PUT("/password", middle.RolePermission(role.USER), api.UserApi.ChangePassword)
			userGroup.DELETE("", middle.RolePermission(role.ROOT), api.UserApi.DeleteUser)
			userGroup.GET("", middle.RolePermission(role.ROOT), api.UserApi.GetUserList)
		}

		pushGroup := apiGroup.Group("/push").Use(middle.RolePermission(role.ADMIN))
		{
			pushGroup.GET("/list", api.PushApi.GetPushList)
			pushGroup.GET("", api.PushApi.GetPushById)
			pushGroup.POST("", api.PushApi.AddPushConfig)
			pushGroup.PUT("", api.PushApi.UpdatePushConfig)
			pushGroup.DELETE("", api.PushApi.DeletePushConfig)
		}

		fileGroup := apiGroup.Group("/file").Use(middle.RolePermission(role.ADMIN))
		{
			fileGroup.GET("/list", api.FileApi.FilePathHandler)
			fileGroup.PUT("", api.FileApi.FileWriteHandler)
			fileGroup.GET("", api.FileApi.FileReadHandler)
		}

		permissionGroup := apiGroup.Group("/permission").Use(middle.RolePermission(role.ROOT))
		{
			permissionGroup.GET("/list", api.PermissionApi.GetPermissionList)
			permissionGroup.PUT("", api.PermissionApi.EditPermssion)
		}

		logGroup := apiGroup.Group("/log").Use(middle.RolePermission(role.ADMIN))
		{
			logGroup.POST("", api.LogApi.GetLog)
		}

		configGroup := apiGroup.Group("/config").Use(middle.RolePermission(role.ROOT))
		{
			configGroup.GET("", api.ConfigApi.GetSystemConfiguration)
			configGroup.PUT("", api.ConfigApi.SetSystemConfiguration)
			configGroup.PUT("/es", api.ConfigApi.EsConfigReload)
		}
	}
}
