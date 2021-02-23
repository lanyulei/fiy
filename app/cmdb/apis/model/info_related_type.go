package model

import (
	"fiy/app/cmdb/models/model"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建模型关系
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

	tx := orm.Eloquent.Begin()

	err = tx.Create(&infoRelatedType).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建关联失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型关系",
		"新建",
		fmt.Sprintf("新建模型关系"),
		map[string]interface{}{},
		infoRelatedType)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 编辑模型关系
func EditInfoRelatedType(c *gin.Context) {
	var (
		err                error
		infoRelatedType    model.InfoRelatedType
		oldInfoRelatedType model.InfoRelatedType
		infoRelatedTypeId  string
	)

	infoRelatedTypeId = c.Param("id")

	err = c.ShouldBind(&infoRelatedType)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldInfoRelatedType, infoRelatedTypeId).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型关系失败")
		return
	}

	tx := orm.Eloquent.Begin()
	err = tx.Model(&model.InfoRelatedType{}).
		Where("id = ?", infoRelatedTypeId).
		Updates(map[string]interface{}{
			"target":          infoRelatedType.Target,
			"related_type_id": infoRelatedType.RelatedTypeId,
			"constraint":      infoRelatedType.Constraint,
			"remark":          infoRelatedType.Remark,
		}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "编辑模型关系失败")
		return
	}

	infoRelatedType.Id = oldInfoRelatedType.Id
	infoRelatedType.Source = oldInfoRelatedType.Source

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型关系",
		"编辑",
		"编辑模型关系",
		oldInfoRelatedType,
		infoRelatedType)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 模型关系列表
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

// 删除模型关系
func DeleteInfoRelatedType(c *gin.Context) {
	var (
		err               error
		infoRelatedTypeId string
		oldData           model.InfoRelatedType
	)

	infoRelatedTypeId = c.Param("id")

	err = orm.Eloquent.Find(&oldData, infoRelatedTypeId).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型关系失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&model.InfoRelatedType{}, infoRelatedTypeId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除模型关系失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型关系",
		"删除",
		fmt.Sprintf("删除模型关系"),
		oldData,
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 模型关系数据
func RelatedInfo(c *gin.Context) {
	var (
		err       error
		modelList []struct {
			Name string `json:"name"`
		}
		relatedList []struct {
			Source string `json:"source"`
			Target string `json:"target"`
		}
	)

	err = orm.Eloquent.Model(&model.Info{}).Select("name").Scan(&modelList).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型列表失败")
		return
	}

	err = orm.Eloquent.
		Model(&model.InfoRelatedType{}).
		Joins("left join cmdb_model_info as s on s.id = cmdb_model_info_related_type.source").
		Joins("left join cmdb_model_info as t on t.id = cmdb_model_info_related_type.target").
		Select("s.name as source, t.name as target").
		Scan(&relatedList).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型关系失败")
		return
	}

	app.OK(c, map[string]interface{}{
		"models":   modelList,
		"relateds": relatedList,
	}, "")
}
