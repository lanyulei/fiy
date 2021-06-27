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
		r.POST("/service-classify", business.CreateServiceClassify)        // 新建
		r.GET("/service-classify", business.ServiceClassifyList)           // 列表
		r.DELETE("/service-classify/:id", business.DeleteServiceClassify)  // 删除
		r.PUT("/service-classify/:id", business.EditServiceClassify)       // 编辑
		r.GET("/svc-tpl", business.ServiceTemplateList)                    // 列表
		r.POST("/svc-tpl", business.CreateServiceTemplate)                 // 新建
		r.GET("/svc-tpl/:id", business.ServiceTemplateDetails)             // 详情
		r.DELETE("/svc-tpl/:id", business.DeleteServiceTemplate)           // 删除
		r.POST("/svc-tpl-process", business.CreateProcess)                 // 新建进程
		r.PUT("/svc-tpl-process/:id", business.EditProcess)                // 编辑进程
		r.DELETE("/svc-tpl-process/:id", business.DeleteProcess)           // 删除进程
		r.POST("/cluster-tpl", business.CreateClusterTpl)                  // 新建集群模板
		r.GET("/cluster-tpl", business.ClusterTplList)                     // 集群模板列表
		r.PUT("/cluster-tpl/:id", business.EditClusterTpl)                 // 编辑集群模板
		r.DELETE("/cluster-tpl/:id", business.DeleteClusterTpl)            // 删除集群模板
		r.GET("/cluster-svc-tpl/:id", business.ClusterTplReSvcTpl)         // 集群模板对应的服务模板列表+
		r.GET("/tree", business.BusinessTree)                              // 集群模板对应的服务模板列表
		r.POST("/add-business-node", business.AddBusinessNode)             // 新建节点
		r.DELETE("/delete-business-node/:id", business.DeleteBusinessNode) // 删除节点
	}
}
