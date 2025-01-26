package route

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/api"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/middle"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.New()
	r.Use(gin.Recovery())
	if !config.CF.UserTui {
		r.Use(gin.Logger())
	}
	routePathInit(r)
	staticInit(r)
	pprofInit(r)
	err := r.Run(config.CF.Listen)
	log.Logger.Fatalw("服务器启动失败", "err", err)
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
	r.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Data(200, "image/x-icon", utils.UnwarpIgnore(f.ReadFile("templates/favicon.ico")))
	})
}

func pprofInit(r *gin.Engine) {
	if config.CF.PprofEnable {
		pprof.Register(r)
		log.Logger.Debug("启用 pprof")
	}
}

func routePathInit(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use(middle.CheckToken())
	apiGroup.Use(middle.PanicMiddle())
	{
		apiGroup.GET("/ws", middle.OprPermission(constants.OPERATION_TERMINAL), api.WsApi.WebsocketHandle)

		processGroup := apiGroup.Group("/process")
		{
			processGroup.DELETE("", middle.OprPermission(constants.OPERATION_STOP), api.ProcApi.KillProcess)
			processGroup.GET("", api.ProcApi.GetProcessList)
			processGroup.GET("/wait", middle.ProcessWaitCond.WaitGetMiddel, api.ProcApi.GetProcessList)
			processGroup.PUT("", middle.OprPermission(constants.OPERATION_START), api.ProcApi.StartProcess)
			processGroup.PUT("/all", api.ProcApi.StartAllProcess)
			processGroup.DELETE("/all", api.ProcApi.KillAllProcess)
			processGroup.GET("/control", middle.RolePermission(constants.ROLE_ROOT), middle.ProcessWaitCond.WaitTriggerMiddel, api.ProcApi.ProcessControl)

			proConfigGroup := processGroup.Group("/config")
			{
				proConfigGroup.POST("", middle.RolePermission(constants.ROLE_ROOT), middle.ProcessWaitCond.WaitTriggerMiddel, bind(api.ProcApi.CreateNewProcess))
				proConfigGroup.DELETE("", middle.RolePermission(constants.ROLE_ROOT), middle.ProcessWaitCond.WaitTriggerMiddel, api.ProcApi.DeleteNewProcess)
				proConfigGroup.PUT("", middle.RolePermission(constants.ROLE_ROOT), bind(api.ProcApi.UpdateProcessConfig))
				proConfigGroup.GET("", middle.RolePermission(constants.ROLE_ADMIN), api.ProcApi.GetProcessConfig)
			}
		}

		taskGroup := apiGroup.Group("/task").Use(middle.RolePermission(constants.ROLE_ADMIN))
		{
			taskGroup.GET("", api.TaskApi.GetTaskById)
			taskGroup.GET("/all", api.TaskApi.GetTaskList)
			taskGroup.GET("/all/wait", middle.TaskWaitCond.WaitGetMiddel, api.TaskApi.GetTaskList)
			taskGroup.POST("", middle.TaskWaitCond.WaitTriggerMiddel, bind(api.TaskApi.CreateTask))
			taskGroup.DELETE("", middle.TaskWaitCond.WaitTriggerMiddel, api.TaskApi.DeleteTaskById)
			taskGroup.PUT("", middle.TaskWaitCond.WaitTriggerMiddel, bind(api.TaskApi.EditTask))
			taskGroup.PUT("/enable", middle.TaskWaitCond.WaitTriggerMiddel, bind(api.TaskApi.EditTaskEnable))
			taskGroup.GET("/start", api.TaskApi.StartTask)
			taskGroup.GET("/stop", api.TaskApi.StopTask)
			taskGroup.POST("/key", api.TaskApi.CreateTaskApiKey)
			taskGroup.GET("/api-key/:key", api.TaskApi.RunTaskByKey)
		}

		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/login", bind(api.UserApi.LoginHandler))
			userGroup.POST("", middle.RolePermission(constants.ROLE_ROOT), bind((api.UserApi.CreateUser)))
			userGroup.PUT("/password", middle.RolePermission(constants.ROLE_USER), bind(api.UserApi.ChangePassword))
			userGroup.DELETE("", middle.RolePermission(constants.ROLE_ROOT), api.UserApi.DeleteUser)
			userGroup.GET("", middle.RolePermission(constants.ROLE_ROOT), api.UserApi.GetUserList)
		}

		pushGroup := apiGroup.Group("/push").Use(middle.RolePermission(constants.ROLE_ADMIN))
		{
			pushGroup.GET("/list", api.PushApi.GetPushList)
			pushGroup.GET("", api.PushApi.GetPushById)
			pushGroup.POST("", bind(api.PushApi.AddPushConfig))
			pushGroup.PUT("", bind(api.PushApi.UpdatePushConfig))
			pushGroup.DELETE("", api.PushApi.DeletePushConfig)
		}

		fileGroup := apiGroup.Group("/file").Use(middle.RolePermission(constants.ROLE_ADMIN))
		{
			fileGroup.GET("/list", api.FileApi.FilePathHandler)
			fileGroup.PUT("", api.FileApi.FileWriteHandler)
			fileGroup.GET("", api.FileApi.FileReadHandler)
		}

		permissionGroup := apiGroup.Group("/permission").Use(middle.RolePermission(constants.ROLE_ROOT))
		{
			permissionGroup.GET("/list", api.PermissionApi.GetPermissionList)
			permissionGroup.PUT("", middle.ProcessWaitCond.WaitTriggerMiddel, bind(api.PermissionApi.EditPermssion))
		}

		logGroup := apiGroup.Group("/log").Use(middle.RolePermission(constants.ROLE_USER))
		{
			logGroup.POST("", bind(api.LogApi.GetLog))
			logGroup.GET("/running", api.LogApi.GetRunningLog)
		}

		configGroup := apiGroup.Group("/config").Use(middle.RolePermission(constants.ROLE_ROOT))
		{
			configGroup.GET("", api.ConfigApi.GetSystemConfiguration)
			configGroup.PUT("", bind(api.ConfigApi.SetSystemConfiguration))
			configGroup.PUT("/es", api.ConfigApi.EsConfigReload)
		}
	}
}

func bind[T any](f func(ctx *gin.Context, req T)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data T
		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			ctx.JSON(500, gin.H{
				"code": -1,
				"msg":  "Invalid parameters!",
			})
			ctx.Abort()
			return
		}
		f(ctx, data)
	}
}
