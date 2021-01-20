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

	// 查询分组是否存在， 分组唯一标识及名称都不存在，才可创建分组
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

func GetModelDetails(c *gin.Context) {
	var (
		err          error
		fieldDetails struct {
			model.Info
			FieldGroups []struct {
				model.FieldGroup
				Fields []model.Fields `json:"fields"`
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

	app.OK(c, fieldDetails, "")
}
