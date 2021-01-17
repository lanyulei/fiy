package router

import (
	"fiy/app/cmdb/apis/model"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterCmdbModelRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/cmdb/model").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		// 模型分组
		r.GET("/group", model.GetModelList)       // 模型分组列表
		r.POST("/group", model.CreateGroup)       // 创建模型分组
		r.PUT("/group/:id", model.EditGroup)      // 编辑模型分组
		r.DELETE("/group/:id", model.DeleteGroup) // 删除模型分组

		// 模型管理
		r.POST("/info", model.CreateModelInfo) // 创建模型
	}
}
