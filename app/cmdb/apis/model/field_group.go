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

// 创建模型字段分组
func CreateModelFieldGroup(c *gin.Context) {
	var (
		err        error
		fieldGroup model.FieldGroup
	)

	err = c.ShouldBind(&fieldGroup)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 创建字段分组
	err = tx.Create(&fieldGroup).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "创建字段分组失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"新建",
		fmt.Sprintf("新建字段分组 \"%s\"", fieldGroup.Name),
		map[string]interface{}{},
		fieldGroup)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除模型分组
func DeleteFieldGroup(c *gin.Context) {
	var (
		err            error
		fieldGroupId   string
		fieldCount     int64
		fieldGroupInfo model.FieldGroup
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

	err = orm.Eloquent.Find(&fieldGroupInfo, fieldGroupId).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段分组失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"删除",
		fmt.Sprintf("删除字段分组 \"%s\"", fieldGroupInfo.Name),
		fieldGroupInfo,
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	// 删除字段分组
	err = tx.Delete(&model.FieldGroup{}, fieldGroupId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除字段分组失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 编辑字段分组
func EditFieldGroup(c *gin.Context) {
	var (
		err            error
		fieldGroup     model.FieldGroup
		fieldGroupId   string
		fieldGroupInfo model.FieldGroup
	)

	fieldGroupId = c.Param("id")

	err = c.ShouldBind(&fieldGroup)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&fieldGroupInfo, fieldGroupId).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段分组失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"编辑",
		fmt.Sprintf("编辑字段分组 \"%s\"", fieldGroup.Name),
		fieldGroup,
		fieldGroupInfo)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	err = tx.Model(&fieldGroup).Where("id = ?", fieldGroupId).Updates(map[string]interface{}{
		"name":     fieldGroup.Name,
		"sequence": fieldGroup.Sequence,
		"is_fold":  fieldGroup.IsFold,
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新字段分组失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}
