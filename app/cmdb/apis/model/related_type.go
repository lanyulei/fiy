package model

import (
	"fiy/app/cmdb/models/model"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"
	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建关联类型
func AddAssociationType(c *gin.Context) {
	var (
		err              error
		association      model.RelatedType
		associationCount int64
	)

	err = c.ShouldBind(&association)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 判断唯一标识及名称是否唯一
	err = orm.Eloquent.
		Model(&model.RelatedType{}).
		Where("identifies = ? or name = ?", association.Identifies, association.Name).
		Count(&associationCount).Error
	if err != nil {
		app.Error(c, -1, err, "验证唯一标识或者名称的唯一性失败")
		return
	}
	if associationCount > 0 {
		app.Error(c, -1, nil, "唯一标识或者名称出现重复，请确认")
		return
	}

	// 新建关联类型
	err = orm.Eloquent.Create(&association).Error
	if err != nil {
		app.Error(c, -1, err, "新建关联类型失败")
		return
	}

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
		err                  error
		associationType      model.RelatedType
		associationTypeCount int64
		associationTypeId    string
	)

	associationTypeId = c.Param("id")

	err = c.ShouldBind(&associationType)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 唯一标识及名称不可重复
	err = orm.Eloquent.
		Model(&associationType).
		Where("identifies = ? or name = ?", associationType.Identifies, associationType.Name).
		Count(&associationTypeCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询关联类型是否存在失败")
		return
	}
	if associationTypeCount > 0 {
		app.Error(c, -1, nil, "唯一标识或名称已存在")
		return
	}

	// 更新关联类型
	err = orm.Eloquent.Model(&associationType).
		Where("id = ?", associationTypeId).
		Updates(map[string]interface{}{
			"identifies":      associationType.Identifies,
			"name":            associationType.Name,
			"source_describe": associationType.SourceDescribe,
			"target_describe": associationType.TargetDescribe,
			"direction":       associationType.Direction,
		}).Error
	if err != nil {
		app.Error(c, -1, err, "更新关联类型失败")
		return
	}

	app.OK(c, nil, "")
}

// 删除关联类型
func DeleteAssociationType(c *gin.Context) {
	var (
		err               error
		associationTypeId string
	)

	associationTypeId = c.Param("id")

	err = orm.Eloquent.Delete(&model.RelatedType{}, associationTypeId).Error
	if err != nil {
		app.Error(c, -1, err, "删除关联分类失败")
		return
	}

	app.OK(c, nil, "")
}
