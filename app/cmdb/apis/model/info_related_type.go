package model

import (
	"fiy/app/cmdb/models/model"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建模型关联
func CreateInfoRelatedType(c *gin.Context) {
	var (
		err             error
		infoRelatedType model.InfoRelatedType
	)

	err = c.ShouldBind(&infoRelatedType)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Create(&infoRelatedType).Error
	if err != nil {
		app.Error(c, -1, err, "新建关联失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑模型关联
func EditInfoRelatedType(c *gin.Context) {
	var (
		err               error
		infoRelatedType   model.InfoRelatedType
		infoRelatedTypeId string
	)

	infoRelatedTypeId = c.Param("id")

	err = c.ShouldBind(&infoRelatedType)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&model.InfoRelatedType{}).
		Where("id = ?", infoRelatedTypeId).
		Updates(map[string]interface{}{
			"target":          infoRelatedType.Target,
			"related_type_id": infoRelatedType.RelatedTypeId,
			"constraint":      infoRelatedType.Constraint,
			"remark":          infoRelatedType.Remark,
		}).Error
	if err != nil {
		app.Error(c, -1, err, "编辑模型关联失败")
		return
	}

	app.OK(c, nil, "")
}

// 模型关联列表
func InfoRelatedTypeList(c *gin.Context) {
	var (
		err     error
		modelId string
		list    []struct {
			Id              int    `json:"id"`
			Source          int    `json:"source"`
			Target          int    `json:"target"`
			RelatedTypeId   int    `json:"related_type_id"`
			Constraint      int    `json:"constraint"`
			Remark          string `json:"remark"`
			RelatedTypeName string `json:"related_type_name"`
			SourceModelName string `json:"source_model_name"`
			TargetModelName string `json:"target_model_name"`
		}
	)

	modelId = c.Param("id")

	err = orm.Eloquent.Model(&model.InfoRelatedType{}).
		Joins("left join cmdb_model_info as source_info on source_info.id = cmdb_model_info_related_type.source").
		Joins("left join cmdb_model_info as target_info on target_info.id = cmdb_model_info_related_type.target").
		Joins("left join cmdb_model_related_type as related_type on related_type.id = cmdb_model_info_related_type.related_type_id").
		Select(`cmdb_model_info_related_type.id,
					  cmdb_model_info_related_type.source,
					  cmdb_model_info_related_type.target,
					  cmdb_model_info_related_type.related_type_id,
					  cmdb_model_info_related_type.constraint,
					  cmdb_model_info_related_type.remark,
					  cmdb_model_info_related_type.remark,
					  related_type.name as related_type_name,
					  source_info.name as source_model_name,
					  target_info.name as target_model_name`).
		Where("cmdb_model_info_related_type.source = ?", modelId).
		Scan(&list).
		Error
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	app.OK(c, list, "")
}

// 删除模型关联
func DeleteInfoRelatedType(c *gin.Context) {
	var (
		err               error
		infoRelatedTypeId string
	)

	infoRelatedTypeId = c.Param("id")

	err = orm.Eloquent.Delete(&model.InfoRelatedType{}, infoRelatedTypeId).Error
	if err != nil {
		app.Error(c, -1, err, "删除模型关联失败")
		return
	}

	app.OK(c, nil, "")
}
