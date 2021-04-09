package model

import (
	"encoding/json"
	"fiy/app/cmdb/models/model"
	"fiy/app/cmdb/models/resource"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 创建模型
func CreateModelInfo(c *gin.Context) {
	var (
		err  error
		info model.Info
	)

	err = c.ShouldBind(&info)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	info.IsUsable = true

	tx := orm.Eloquent.Begin()

	// 写入数据库
	err = tx.Create(&info).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "创建模型失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"新建",
		fmt.Sprintf("新建模型 \"%s\"", info.Name),
		map[string]interface{}{},
		info)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, info, "")
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

	// 查询字段分组
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

// 获取模型对应的所有字段列表
func GetModelFields(c *gin.Context) {
	var (
		err     error
		fields  []model.Fields
		modelId string
	)

	modelId = c.Param("id")

	err = orm.Eloquent.Model(&model.Fields{}).
		Where("info_id = ?", modelId).
		Order("list_display_sort").
		Find(&fields).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段列表失败")
		return
	}

	app.OK(c, fields, "")
}

// 查询当前模型关联的所有模型及对应的字段
func GetRelatedModelFields(c *gin.Context) {
	var (
		err                error
		modelID            string
		dataID             string
		relatedModelIDs    []int
		relatedModelFields []*struct {
			model.Info
			Fields []model.Fields `json:"fields"`
		}
		dataMap map[string]interface{}
	)

	modelID = c.Param("id")
	dataID = c.DefaultQuery("data_id", "")
	if dataID == "" {
		app.Error(c, -1, nil, "参数data_id不能为空")
		return
	}

	// 查询关联的模型ID
	err = orm.Eloquent.Model(&model.InfoRelatedType{}).
		Where("source = ?", modelID).
		Pluck("target", &relatedModelIDs).Error
	if err != nil {
		app.Error(c, -1, err, "查询关联的模型ID失败")
		return
	}

	// 查询关联模型的信息
	err = orm.Eloquent.Model(&model.Info{}).
		Where("id in ?", relatedModelIDs).
		Find(&relatedModelFields).Error
	if err != nil {
		app.Error(c, -1, err, "查询关联模型信息失败")
		return
	}

	// 查询关联模型的字段
	dataMap = make(map[string]interface{})
	for _, f := range relatedModelFields {
		err = orm.Eloquent.Model(&model.Fields{}).
			Where("info_id = ?", f.Id).
			Find(&f.Fields).Error
		if err != nil {
			app.Error(c, -1, err, "查询关联模型的字段失败")
			return
		}

		var (
			datas []struct {
				resource.Data
				RelatedID int `json:"related_id"`
			}
		)
		dataList := make([]map[string]interface{}, 0)
		err = orm.Eloquent.Model(&resource.DataRelated{}).
			Joins("left join cmdb_resource_data as d on cmdb_resource_data_related.target = d.id").
			Select("d.*, cmdb_resource_data_related.id as related_id").
			Where("cmdb_resource_data_related.source = ? and cmdb_resource_data_related.target_info_id = ?", dataID, f.Id).
			Find(&datas).Error
		if err != nil {
			app.Error(c, -1, err, "查询关联数据ID失败")
			return
		}

		for _, d := range datas {
			var dataJsonMap map[string]interface{}
			err = json.Unmarshal(d.Data.Data, &dataJsonMap)
			if err != nil {
				app.Error(c, -1, err, "反序列化失败")
				return
			}
			dataJsonMap["id"] = d.Id                // 数据ID
			dataJsonMap["info_id"] = d.InfoId       // 模型ID
			dataJsonMap["related_id"] = d.RelatedID // 关联ID
			dataList = append(dataList, dataJsonMap)
		}

		dataMap[f.Identifies] = dataList
	}

	app.OK(c, map[string]interface{}{
		"fields": relatedModelFields,
		"data":   dataMap,
	}, "")
}

// 编辑模型
func EditModelInfo(c *gin.Context) {
	var (
		err     error
		info    model.Info
		oldInfo model.Info
		infoId  string
	)

	infoId = c.Param("id")

	err = c.ShouldBind(&info)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldInfo, infoId).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型数据失败")
		return
	}

	tx := orm.Eloquent.Begin()
	newData := map[string]interface{}{
		"identifies": info.Identifies,
		"name":       info.Name,
		"icon":       info.Icon,
		"group_id":   info.GroupId,
	}
	err = tx.Model(&info).Where("id = ?", infoId).Updates(newData).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新模型数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"编辑",
		fmt.Sprintf("编辑模型 \"%s\"", info.Name),
		oldInfo,
		newData)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 停用模型
func StopModelInfo(c *gin.Context) {
	var (
		err     error
		modelId string
		oldData model.Info
		params  struct {
			IsUsable bool `json:"is_usable"`
		}
		memo string
	)

	modelId = c.Param("id")

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, modelId).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Model(&model.Info{}).
		Where("id = ?", modelId).
		Update("is_usable", params.IsUsable).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新模型状态")
		return
	}

	if params.IsUsable {
		memo = "开启"
	} else {
		memo = "停用"
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"编辑",
		fmt.Sprintf("%s模型 \"%s\"", memo, oldData.Name),
		map[string]interface{}{},
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

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
		err           error
		fieldId       string
		uniqueStatus  string
		isUnique      bool
		fieldValue    model.Fields
		oldFieldValue model.Fields
		memo          string
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

	err = orm.Eloquent.Find(&oldFieldValue, fieldId).Error
	if err != nil {
		app.Error(c, -1, err, "查询字段数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Model(&model.Fields{}).
		Where("id = ?", fieldId).
		Update("is_unique", isUnique).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新唯一校验失败")
		return
	}

	if isUnique {
		memo = "新建"
	} else {
		memo = "删除"
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		memo,
		fmt.Sprintf("模型ID：%d, %s字段唯一校验 \"%s\"", oldFieldValue.InfoId, memo, oldFieldValue.Name),
		map[string]interface{}{},
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除模型
func DeleteModelInfo(c *gin.Context) {
	var (
		err     error
		modelId string
		oldInfo model.Info
	)

	modelId = c.Param("id")

	err = orm.Eloquent.Find(&oldInfo, modelId).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&model.Info{}, modelId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除模型失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"删除",
		fmt.Sprintf("删除模型 \"%s\"", oldInfo.Name),
		oldInfo,
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}
