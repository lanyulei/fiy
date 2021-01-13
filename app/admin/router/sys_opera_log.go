package router

import (
	"fiy/app/admin/apis/system/sys_opera_log"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

// 需认证的路由代码
func RegisterSysOperaLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &sys_opera_log.SysOperaLog{}
	r := v1.Group("/sys-opera-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysOperaLogList)
		r.GET("/:id", api.GetSysOperaLog)
		r.POST("", api.InsertSysOperaLog)
		r.PUT("/:id", api.UpdateSysOperaLog)
		r.DELETE("", api.DeleteSysOperaLog)
	}
}
