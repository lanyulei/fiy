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

// 创建模型字段分组
func CreateModelFieldGroup(c *gin.Context) {
	var (
		err             error
		fieldGroup      model.FieldGroup
		fieldGroupCount int64
	)

	err = c.ShouldBind(&fieldGroup)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 验证字段分组是否存在
	err = orm.Eloquent.Model(&fieldGroup).Where("name = ?", fieldGroup.Name).Count(&fieldGroupCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段分组是否存在失败")
		return
	}
	if fieldGroupCount > 0 {
		app.Error(c, -1, nil, "字段分组名称已存在，请确认")
		return
	}

	// 创建字段分组
	err = orm.Eloquent.Create(&fieldGroup).Error
	if err != nil {
		app.Error(c, -1, err, "创建字段分组失败")
		return
	}

	app.OK(c, nil, "")
}

// 删除模型分组
func DeleteFieldGroup(c *gin.Context) {
	var (
		err          error
		fieldGroupId string
		fieldCount   int64
	)

	fieldGroupId = c.Param("id")

	// 如果分组下有对应字段，则无法删除
	err = orm.Eloquent.Model(&model.Fields{}).Where("field_group_id = ?", fieldGroupId).Count(&fieldCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段列表失败")
		return
	}
	if fieldCount > 0 {
		app.Error(c, -1, err, "无法删除分组，因分组下有对应的字段数据")
		return
	}

	// 删除字段分组
	err = orm.Eloquent.Delete(&model.FieldGroup{}, fieldGroupId).Error
	if err != nil {
		app.Error(c, -1, err, "删除字段分组失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑字段分组
func EditFieldGroup(c *gin.Context) {
	var (
		err             error
		fieldGroup      model.FieldGroup
		fieldGroupId    string
		fieldGroupCount int64
	)

	fieldGroupId = c.Param("id")

	err = c.ShouldBind(&fieldGroup)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 验证字段分组是否存在
	err = orm.Eloquent.Model(&fieldGroup).Where("name = ?", fieldGroup.Name).Count(&fieldGroupCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段分组是否存在失败")
		return
	}
	if fieldGroupCount > 0 {
		app.Error(c, -1, nil, "字段分组名称已存在，请确认")
		return
	}

	err = orm.Eloquent.Model(&fieldGroup).Where("id = ?", fieldGroupId).Updates(map[string]interface{}{
		"name":     fieldGroup.Name,
		"sequence": fieldGroup.Sequence,
		"is_fold":  fieldGroup.IsFold,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "更新字段分组失败")
		return
	}

	app.OK(c, nil, "")
}
