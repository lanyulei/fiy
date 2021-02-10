package router

import (
	"fiy/app/cmdb/apis/business"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterCmdbBusinessRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/cmdb/business").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.POST("/service-classify", business.CreateServiceClassify)       // 新建
		r.GET("/service-classify", business.ServiceClassifyList)          // 列表
		r.DELETE("/service-classify/:id", business.DeleteServiceClassify) // 删除
		r.PUT("/service-classify/:id", business.EditServiceClassify)      // 删除
	}
}
