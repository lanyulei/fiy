package router

import (
	"fiy/app/k8s/apis/version"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterK8sVersionRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/k8s/version").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", version.List)
	}
}
