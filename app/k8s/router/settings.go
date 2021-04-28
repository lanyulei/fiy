package router

import (
	"fiy/app/k8s/apis/settings"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterK8sSettingsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/k8s/settings").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		// 仓库信息
		r.POST("/registry", settings.CreateRegistry)
		r.PUT("/registry/:id", settings.EditRegistry)
		r.GET("/registry", settings.RegistryList)
		r.DELETE("/registry/:id", settings.DeleteRegistry)

		// 凭据
		r.POST("/credential", settings.CreateCredential)
		r.PUT("/credential/:id", settings.EditCredential)
		r.GET("/credential", settings.CredentialList)
		r.DELETE("/credential/:id", settings.DeleteCredential)

		// NTP服务
		r.POST("/ntp", settings.SaveNTP)
		r.GET("/ntp", settings.GetNTP)
	}
}
