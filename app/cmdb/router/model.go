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
		r.GET("/group", model.GetModelList)
		r.POST("/group", model.CreateGroup)
		r.PUT("/group/:id", model.EditGroup)
		r.DELETE("/group/:id", model.DeleteGroup)
		r.GET("/group-list", model.GetModelGroupList)

		// 模型管理
		r.POST("/info", model.CreateModelInfo)
		r.PUT("/info/:id", model.EditModelInfo)
		r.DELETE("/info/:id", model.DeleteModelInfo)
		r.PUT("/stop/info/:id", model.StopModelInfo)

		// 模型详情
		r.POST("/field-group", model.CreateModelFieldGroup)
		r.GET("/details/:id", model.GetModelDetails)
		r.GET("/fields/:id", model.GetModelFields)
		r.GET("/related-fields/:id", model.GetRelatedModelFields)
		r.POST("/field", model.CreateModelField)
		r.PUT("/field/:id", model.EditModelField)
		r.DELETE("/field/:id", model.DeleteModelField)
		r.DELETE("/field-group/:id", model.DeleteFieldGroup)
		r.PUT("/field-group/:id", model.EditFieldGroup)
		r.GET("/unique-field/:id", model.GetModelUniqueFields)
		r.PUT("/unique-field/:id", model.UpdateFieldUnique)
		r.POST("/model-related", model.CreateInfoRelatedType)
		r.GET("/model-related/:id", model.InfoRelatedTypeList)
		r.PUT("/model-related/:id", model.EditInfoRelatedType)
		r.DELETE("/model-related/:id", model.DeleteInfoRelatedType)

		// 关联类型
		r.POST("/association-type", model.AddAssociationType)
		r.GET("/association-type", model.AssociationTypeList)
		r.PUT("/association-type/:id", model.UpdateAssociationType)
		r.DELETE("/association-type/:id", model.DeleteAssociationType)

		// 模型关系
		r.GET("/related", model.RelatedInfo)
	}
}
