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

// 创建模型字段
func CreateModelField(c *gin.Context) {
	var (
		err        error
		fieldValue model.Fields
	)

	err = c.ShouldBind(&fieldValue)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	if fieldValue.Identifies == "id" ||
		fieldValue.Identifies == "uuid" ||
		fieldValue.Identifies == "info_id" ||
		fieldValue.Identifies == "related_id" {
		app.Error(c, -1, err, "id、info_id、uuid 和 related_id 是预留的字段标识，请选择其他的字段标识")
		return
	}

	tx := orm.Eloquent.Begin()

	// 创建字段
	err = tx.Create(&fieldValue).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "创建字段失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"新建",
		fmt.Sprintf("新建字段 \"%s\"", fieldValue.Name),
		map[string]interface{}{},
		fieldValue)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 更新模型字段
func EditModelField(c *gin.Context) {
	var (
		err     error
		field   model.Fields
		fieldId string
		oldData model.Fields
	)

	fieldId = c.Param("id")

	err = c.ShouldBind(&field)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, fieldId).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段信息失败")
		return
	}

	tx := orm.Eloquent.Begin()

	newData := map[string]interface{}{
		"identifies":        field.Identifies,
		"name":              field.Name,
		"type":              field.Type,
		"is_edit":           field.IsEdit,
		"is_unique":         field.IsUnique,
		"required":          field.Required,
		"prompt":            field.Prompt,
		"configuration":     field.Configuration,
		"is_list_display":   field.IsListDisplay,
		"list_display_sort": field.ListDisplaySort,
	}
	err = tx.Model(&field).Where("id = ?", fieldId).Updates(newData).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"编辑",
		fmt.Sprintf("编辑字段 \"%s\"", field.Name),
		oldData,
		newData)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除模型字段
func DeleteModelField(c *gin.Context) {
	var (
		err     error
		fieldId string
		oldData model.Fields
	)

	fieldId = c.Param("id")

	err = orm.Eloquent.Find(&oldData, fieldId).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段信息失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&model.Fields{}, fieldId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除模型字段失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"删除",
		fmt.Sprintf("删除字段 \"%s\"", oldData.Name),
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
