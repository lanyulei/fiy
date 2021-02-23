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
  @Desc : 模型分组
*/

// 创建模型分组
func CreateGroup(c *gin.Context) {
	var (
		err   error
		group model.Group
	)

	err = c.ShouldBind(&group)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 写入数据库
	err = tx.Create(&group).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "创建分组失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"新建",
		fmt.Sprintf("新建模型分组 <%s>", group.Name),
		map[string]interface{}{},
		group)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

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
		oldData model.Group
		groupId string
	)

	groupId = c.Param("id")

	err = c.ShouldBind(&group)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, groupId).Error
	if err != nil {
		app.Error(c, -1, err, "获取模型分组失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Model(&group).Where("id = ?", groupId).Updates(group).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "编辑模型分组失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"编辑",
		fmt.Sprintf("编辑模型分组 <%s>", group.Name),
		oldData,
		group)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除模型分组
func DeleteGroup(c *gin.Context) {
	var (
		err     error
		groupId string
		oldData model.Group
	)

	groupId = c.Param("id")

	err = orm.Eloquent.Find(&oldData, groupId).Error
	if err != nil {
		app.Error(c, -1, err, "获取模型分组失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&model.Group{}, groupId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "模型分组删除失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"模型",
		"模型管理",
		"删除",
		fmt.Sprintf("删除模型分组 <%s>", oldData.Name),
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

// 获取模型分组列表
func GetModelGroupList(c *gin.Context) {
	var (
		err       error
		groupList []model.Group
	)

	err = orm.Eloquent.Model(&model.Group{}).Find(&groupList).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型分组失败")
		return
	}

	app.OK(c, groupList, "")
}
