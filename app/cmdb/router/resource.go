package router

import (
	"fiy/app/cmdb/apis/resource"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterCmdbResourceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/cmdb/resource").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.POST("", resource.CreateData) // 新建资源数据
	}
}
