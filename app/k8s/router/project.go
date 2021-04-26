package router

import (
	"fiy/app/k8s/apis/project"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterK8sProjectRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/k8s/project").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", project.List)          // 项目列表
		r.POST("", project.Create)       // 创建项目
		r.PUT("/:id", project.Edit)      // 编辑项目
		r.DELETE("/:id", project.Delete) // 删除项目
	}
}
