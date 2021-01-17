package model

import (
	"fiy/app/cmdb/models/model"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
  @Desc : 模型分组
*/

// 创建模型分组
func CreateGroup(c *gin.Context) {
	var (
		err        error
		group      model.Group
		groupCount int64
	)

	err = c.ShouldBind(&group)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 查询分组是否存在， 分组唯一标识及名称都不存在，才可创建分组
	err = orm.Eloquent.
		Model(&group).
		Where("identifies = ? or name = ?", group.Identifies, group.Name).
		Count(&groupCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询分组是否存在失败")
		return
	}
	if groupCount > 0 {
		app.Error(c, -1, nil, "分组唯一标识或名称已存在")
		return
	}

	// 写入数据库
	err = orm.Eloquent.Create(&group).Error
	if err != nil {
		app.Error(c, -1, err, "创建分组失败")
		return
	}

	app.OK(c, nil, "")
}

// 模型列表
func GetModelList(c *gin.Context) {
	var (
		err       error
		modelList []*struct {
			model.Group
			ModelList []model.Info `json:"model_list"`
		}
	)

	isUsable := c.DefaultQuery("isUsable", "1")
	search := c.DefaultQuery("search", "")

	// 获取分组
	err = orm.Eloquent.Model(&model.Group{}).Find(&modelList).Error
	if err != nil {
		app.Error(c, -1, err, "获取分组列表失败")
		return
	}

	// 获取分组对应的模型列表
	for _, group := range modelList {
		db := orm.Eloquent.Model(&model.Info{})
		if search != "" {
			db = db.Where("name like ?", "%"+search+"%")
		}
		err = db.Where("group_id = ? and is_usable = ?", group.Id, isUsable).
			Find(&group.ModelList).Error
		if err != nil {
			app.Error(c, -1, err, "获取模型信息失败")
			return
		}
	}

	app.OK(c, modelList, "")
}

// 编辑模型分组
func EditGroup(c *gin.Context) {
	var (
		err     error
		group   model.Group
		groupId string
	)

	groupId = c.Param("id")

	err = c.ShouldBind(&group)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&group).Where("id = ?", groupId).Updates(group).Error
	if err != nil {
		app.Error(c, -1, err, "编辑模型分组失败")
		return
	}

	app.OK(c, nil, "")
}

// 删除模型分组
func DeleteGroup(c *gin.Context) {
	var (
		err     error
		groupId string
	)

	groupId = c.Param("id")

	err = orm.Eloquent.Delete(&model.Group{}, groupId).Error
	if err != nil {
		app.Error(c, -1, err, "模型分组删除失败")
		return
	}

	app.OK(c, nil, "")
}
