package model

import (
	"fiy/app/cmdb/models/model"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建关联类型
func AddAssociationType(c *gin.Context) {
	var (
		err         error
		association model.RelatedType
	)

	err = c.ShouldBind(&association)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 新建关联类型
	err = tx.Create(&association).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建关联类型失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"关联类型",
		"新建",
		fmt.Sprintf("新建关联类型 \"%s\"", association.Name),
		map[string]interface{}{},
		association)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 关联类型列表
func AssociationTypeList(c *gin.Context) {
	var (
		err             error
		result          interface{}
		associationList []*model.RelatedType
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: orm.Eloquent.Model(&model.RelatedType{}),
	}, &associationList, SearchParams)
	if err != nil {
		app.Error(c, -1, err, "分页查询关联类型列表失败")
		return
	}

	app.OK(c, result, "")
}

// 编辑关联类型
func UpdateAssociationType(c *gin.Context) {
	var (
		err               error
		associationType   model.RelatedType
		associationTypeId string
		oldData           model.RelatedType
	)

	associationTypeId = c.Param("id")

	err = c.ShouldBind(&associationType)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, associationTypeId).Error
	if err != nil {
		app.Error(c, -1, err, "查询关联类型失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 更新关联类型
	err = tx.Model(&associationType).
		Where("id = ?", associationTypeId).
		Updates(map[string]interface{}{
			"identifies":      associationType.Identifies,
			"name":            associationType.Name,
			"source_describe": associationType.SourceDescribe,
			"target_describe": associationType.TargetDescribe,
			"direction":       associationType.Direction,
		}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新关联类型失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"关联类型",
		"编辑",
		fmt.Sprintf("编辑关联类型 \"%s\"", associationType.Name),
		oldData,
		associationType)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除关联类型
func DeleteAssociationType(c *gin.Context) {
	var (
		err               error
		associationTypeId string
		oldData           model.RelatedType
	)

	associationTypeId = c.Param("id")

	err = orm.Eloquent.Find(&oldData, associationTypeId).Error
	if err != nil {
		app.Error(c, -1, err, "查询关联类型失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&model.RelatedType{}, associationTypeId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除关联分类失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"关联类型",
		"删除",
		fmt.Sprintf("删除关联类型 \"%s\"", oldData.Name),
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
