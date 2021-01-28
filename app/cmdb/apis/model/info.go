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

// 创建模型
func CreateModelInfo(c *gin.Context) {
	var (
		err       error
		info      model.Info
		infoCount int64
	)

	err = c.ShouldBind(&info)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 唯一标识及名称不可重复
	info.IsUsable = true
	err = orm.Eloquent.
		Model(&info).
		Where("identifies = ? or name = ?", info.Identifies, info.Name).
		Count(&infoCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型是否存在失败")
		return
	}
	if infoCount > 0 {
		app.Error(c, -1, nil, "模型唯一标识或名称已存在")
		return
	}

	// 写入数据库
	err = orm.Eloquent.Create(&info).Error
	if err != nil {
		app.Error(c, -1, err, "创建模型失败")
		return
	}

	app.OK(c, nil, "")
}

// 获取模型详情
func GetModelDetails(c *gin.Context) {
	var (
		err          error
		fieldDetails struct {
			model.Info
			FieldGroups []*struct {
				model.FieldGroup
				Fields []*model.Fields `json:"fields"`
			} `json:"field_groups"`
		}
		modelId string
	)

	modelId = c.Param("id")

	// 查询模型信息
	err = orm.Eloquent.Model(&model.Info{}).Where("id = ?", modelId).Find(&fieldDetails).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型信息失败")
		return
	}

	// 查询模型分组
	err = orm.Eloquent.Model(&model.FieldGroup{}).Where("info_id = ?", modelId).Find(&fieldDetails.FieldGroups).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型信息失败")
		return
	}

	// 获取分组对应的字段
	for _, group := range fieldDetails.FieldGroups {
		err = orm.Eloquent.Model(&model.Fields{}).
			Where("info_id = ? and field_group_id = ?", modelId, group.Id).
			Find(&group.Fields).Error
		if err != nil {
			app.Error(c, -1, err, "查询字段列表失败")
			return
		}
	}

	app.OK(c, fieldDetails, "")
}

// 编辑模型
func EditModelInfo(c *gin.Context) {
	var (
		err       error
		info      model.Info
		infoId    string
		infoCount int64
	)

	infoId = c.Param("id")

	err = c.ShouldBind(&info)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 唯一标识及名称不可重复
	info.IsUsable = true
	err = orm.Eloquent.
		Model(&info).
		Where("identifies = ? or name = ?", info.Identifies, info.Name).
		Count(&infoCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型是否存在失败")
		return
	}
	if infoCount > 0 {
		app.Error(c, -1, nil, "模型唯一标识或名称已存在")
		return
	}

	err = orm.Eloquent.Model(&info).Where("id = ?", infoId).Updates(map[string]interface{}{
		"identifies": info.Identifies,
		"name":       info.Name,
		"icon":       info.Icon,
		"group_id":   info.GroupId,
	}).Error

	if err != nil {
		app.Error(c, -1, err, "更新模型数据失败")
		return
	}

	app.OK(c, nil, "")
}

// 停用模型
func StopModelInfo(c *gin.Context) {
	var (
		err     error
		modelId string
		params  struct {
			IsUsable bool `json:"is_usable"`
		}
	)

	modelId = c.Param("id")

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&model.Info{}).
		Where("id = ?", modelId).
		Update("is_usable", params.IsUsable).Error
	if err != nil {
		app.Error(c, -1, err, "更新模型状态")
		return
	}

	app.OK(c, nil, "")
}

// 获模型中唯一校验的列
func GetModelUniqueFields(c *gin.Context) {
	var (
		err     error
		modelId string
		fields  []model.Fields
	)

	modelId = c.Param("id")

	err = orm.Eloquent.Model(&model.Fields{}).
		Select("id, identifies, name").
		Where("info_id = ? and is_unique = 1", modelId).
		Find(&fields).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段数据失败")
		return
	}

	app.OK(c, fields, "")
}

// 更新字段唯一校验规则
func UpdateFieldUnique(c *gin.Context) {
	var (
		err          error
		fieldId      string
		uniqueStatus string
		isUnique     bool
		fieldValue   model.Fields
	)

	fieldId = c.Param("id")

	uniqueStatus = c.DefaultQuery("unique_status", "")
	if uniqueStatus == "" {
		app.Error(c, -1, err, "unique_status 参数异常请确认")
		return
	}

	if uniqueStatus == "create" {
		// 校验是否已经开启唯一校验
		err = orm.Eloquent.Model(&model.Fields{}).
			Select("id, is_unique").
			Where("id = ?", fieldId).
			Find(&fieldValue).Error
		if err != nil {
			app.Error(c, -1, err, "查新字段唯一校验状态失败")
			return
		}
		if fieldValue.IsUnique {
			app.Error(c, -1, err, "相同的唯一校验规则已经存在")
			return
		}
		isUnique = true
	} else if uniqueStatus == "delete" {
		isUnique = false
	}

	err = orm.Eloquent.
		Model(&model.Fields{}).
		Where("id = ?", fieldId).
		Update("is_unique", isUnique).Error
	if err != nil {
		app.Error(c, -1, err, "更新唯一校验失败")
		return
	}

	app.OK(c, nil, "")
}
