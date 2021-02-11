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
		r.PUT("/service-classify/:id", business.EditServiceClassify)      // 编辑

		r.GET("/svc-tpl", business.ServiceTemplateList)        // 列表
		r.POST("/svc-tpl", business.CreateServiceTemplate)     // 新建
		r.GET("/svc-tpl/:id", business.ServiceTemplateDetails) // 详情

		r.POST("/svc-tpl-process", business.CreateProcess)  // 新建流程
		r.PUT("/svc-tpl-process/:id", business.EditProcess) // 编辑流程
	}
}
