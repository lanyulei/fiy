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

// 创建模型字段
func CreateModelField(c *gin.Context) {
	var (
		err        error
		fieldValue model.Fields
		fieldCount int64
	)

	err = c.ShouldBind(&fieldValue)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 判断唯一标识及名称是否唯一
	err = orm.Eloquent.
		Model(&model.Fields{}).
		Where("info_id = ? and (identifies = ? or name = ?)", fieldValue.InfoId, fieldValue.Identifies, fieldValue.Name).
		Count(&fieldCount).Error
	if err != nil {
		app.Error(c, -1, err, "验证唯一标识或者名称的唯一性失败")
		return
	}
	if fieldCount > 0 {
		app.Error(c, -1, nil, "唯一标识或者名称出现重复，请确认")
		return
	}

	// 创建字段
	err = orm.Eloquent.Create(&fieldValue).Error
	if err != nil {
		app.Error(c, -1, err, "创建字段失败")
		return
	}

	app.OK(c, nil, "")
}

// 更新模型字段
func EditModelField(c *gin.Context) {
	var (
		err        error
		field      model.Fields
		fieldId    string
		fieldCount int64
	)

	fieldId = c.Param("id")

	err = c.ShouldBind(&field)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 判断唯一标识及名称是否唯一
	err = orm.Eloquent.
		Model(&model.Fields{}).
		Where("info_id = ? and (identifies = ? or name = ?)", field.InfoId, field.Identifies, field.Name).
		Count(&fieldCount).Error
	if err != nil {
		app.Error(c, -1, err, "验证唯一标识或者名称的唯一性失败")
		return
	}
	if fieldCount > 0 {
		app.Error(c, -1, nil, "唯一标识或者名称出现重复，请确认")
		return
	}

	err = orm.Eloquent.Model(&field).Where("id = ?", fieldId).Updates(map[string]interface{}{
		"identifies":    field.Identifies,
		"name":          field.Name,
		"type":          field.Type,
		"is_edit":       field.IsEdit,
		"is_unique":     field.IsUnique,
		"required":      field.Required,
		"prompt":        field.Prompt,
		"configuration": field.Configuration,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	app.OK(c, nil, "")
}

// 删除模型字段
func DeleteModelField(c *gin.Context) {
	var (
		err     error
		fieldId string
	)

	fieldId = c.Param("id")

	err = orm.Eloquent.Delete(&model.Fields{}, fieldId).Error
	if err != nil {
		app.Error(c, -1, err, "删除模型字段失败")
		return
	}

	app.OK(c, nil, "")
}
