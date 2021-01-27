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
		app.Error(c, -1, nil, "新建关联类型失败")
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

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: orm.Eloquent.Model(&model.RelatedType{}),
	}, &associationList)
	if err != nil {
		app.Error(c, -1, err, "分页查询关联类型列表失败")
		return
	}

	app.OK(c, result, "")
}
