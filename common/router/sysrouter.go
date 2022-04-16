package router

import (
	"fiy/app/admin/apis/monitor"
	"fiy/app/admin/apis/system"
	"fiy/app/admin/apis/system/dict"
	. "fiy/app/admin/apis/tools"
	adminRouter "fiy/app/admin/router"
	cmdbRouter "fiy/app/cmdb/router"
	"fiy/common/middleware/handler"
	_ "fiy/docs"
	jwt "fiy/pkg/jwtauth"
	"fiy/pkg/ws"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	sysBaseRouter(g)
	// 静态文件
	sysStaticFileRouter(g, r)
	// swagger；注意：生产环境可以注释掉
	sysSwaggerRouter(g)
	// 无需认证
	sysNoCheckRoleRouter(g)
	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)
	return g
}

func sysBaseRouter(r *gin.RouterGroup) {

	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	r.GET("/", system.Index)
	r.GET("/info", handler.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup, g *gin.Engine) {
	r.Static("/static", "./static/ui/static")
	g.LoadHTMLGlob("static/ui/index.html")
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func sysNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")

	v1.GET("/monitor/server", monitor.ServerInfo)
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)
	v1.GET("/gen/preview/:tableId", Preview)
	v1.GET("/gen/toproject/:tableId", GenCodeV3)
	v1.GET("/gen/todb/:tableId", GenMenuAndApi)
	v1.GET("/gen/tabletree", GetSysTablesTree)
	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
	v1.GET("/dict/databytype/:dictType", dict.GetDictDataByDictType)

	registerDBRouter(v1)
	registerSysTableRouter(v1)
	adminRouter.RegisterPublicRouter(v1)
	adminRouter.RegisterSysSettingRouter(v1)
}

func registerDBRouter(api *gin.RouterGroup) {
	db := api.Group("/db")
	{
		db.GET("/tables/page", GetDBTableList)
		db.GET("/columns/page", GetDBColumnList)
	}
}

func registerSysTableRouter(v1 *gin.RouterGroup) {
	systables := v1.Group("/sys/tables")
	{
		systables.GET("/page", GetSysTableList)
		tablesinfo := systables.Group("/info")
		{
			tablesinfo.POST("", InsertSysTable)
			tablesinfo.PUT("", UpdateSysTable)
			tablesinfo.DELETE("/:tableId", DeleteSysTables)
			tablesinfo.GET("/:tableId", GetSysTables)
			tablesinfo.GET("", GetSysTablesInfo)
		}
	}
}

func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r.POST("/login", authMiddleware.LoginHandler)
	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.Group("").Use(authMiddleware.MiddlewareFunc()).GET("/ws/:id/:channel", ws.WebsocketManager.WsClient)
	r.Group("").Use(authMiddleware.MiddlewareFunc()).GET("/wslogout/:id/:channel", ws.WebsocketManager.UnWsClient)
	v1 := r.Group("/api/v1")

	// admin
	adminRouter.RegisterPageRouter(v1, authMiddleware)
	adminRouter.RegisterBaseRouter(v1, authMiddleware)
	adminRouter.RegisterDeptRouter(v1, authMiddleware)
	adminRouter.RegisterDictRouter(v1, authMiddleware)
	adminRouter.RegisterSysUserRouter(v1, authMiddleware)
	adminRouter.RegisterRoleRouter(v1, authMiddleware)
	adminRouter.RegisterUserCenterRouter(v1, authMiddleware)
	adminRouter.RegisterPostRouter(v1, authMiddleware)
	adminRouter.RegisterMenuRouter(v1, authMiddleware)
	adminRouter.RegisterSysConfigRouter(v1, authMiddleware)
	adminRouter.RegisterSysLoginLogRouter(v1, authMiddleware)
	adminRouter.RegisterSysOperaLogRouter(v1, authMiddleware)

	// cmdb
	cmdbRouter.RegisterCmdbModelRouter(v1, authMiddleware)
	cmdbRouter.RegisterCmdbResourceRouter(v1, authMiddleware)
	cmdbRouter.RegisterCmdbBusinessRouter(v1, authMiddleware)
	cmdbRouter.RegisterCmdbAnalysisRouter(v1, authMiddleware)
	cmdbRouter.RegisterCmdbSearchRouter(v1, authMiddleware)
}
