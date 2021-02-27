package router

import (
	"fiy/app/cmdb/apis/analysis"
	"fiy/common/middleware"
	jwt "fiy/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterCmdbAnalysisRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/cmdb/analysis").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", analysis.AuditList)                // 操作列表
		r.GET("/details/:id", analysis.AuditDetails) // 详情

		r.GET("/operation", analysis.Operation)
	}
}
