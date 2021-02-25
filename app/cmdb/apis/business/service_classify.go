package business

import (
	"fiy/app/cmdb/models/business"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 创建服务分类
func CreateServiceClassify(c *gin.Context) {
	var (
		err             error
		remark          string
		serviceClassify business.ServiceClassify
	)

	err = c.ShouldBind(&serviceClassify)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 写入数据库
	err = tx.Create(&serviceClassify).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "创建服务分类失败")
		return
	}

	// 添加操作审计
	if serviceClassify.Level == 1 {
		remark = "新建分组名称"
	} else if serviceClassify.Level == 2 {
		remark = "新建服务分类"
	}
	err = actions.AddAudit(c, tx, "业务", "服务分类", "新建", fmt.Sprintf("%s \"%s\"", remark, serviceClassify.Name), map[string]interface{}{}, serviceClassify)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 服务分类列表
func ServiceClassifyList(c *gin.Context) {
	var (
		err  error
		list []*struct {
			business.ServiceClassify
			Services []business.ServiceClassify `json:"services"`
		}
	)

	err = orm.Eloquent.Model(&business.ServiceClassify{}).
		Where("level = ?", 1).
		Find(&list).Error
	if err != nil {
		app.Error(c, -1, err, "查询分组列表失败")
		return
	}

	for _, g := range list {
		err = orm.Eloquent.Model(&business.ServiceClassify{}).
			Where("level = ? and parent = ?", 2, g.Id).
			Find(&g.Services).Error
		if err != nil {
			app.Error(c, -1, err, "查询服务分类列表失败")
			return
		}
	}

	app.OK(c, list, "")
}

// 删除服务分类
func DeleteServiceClassify(c *gin.Context) {
	var (
		err       error
		id        string
		dataCount int64
		level     string
		oldData   business.ServiceClassify
	)

	id = c.Param("id")

	level = c.DefaultQuery("level", "")
	if level == "" {
		app.Error(c, -1, err, "参数level不存在，请确认")
		return
	}

	if level == "1" {
		err = orm.Eloquent.Model(&business.ServiceClassify{}).Where("level = ? and parent = ?", 2, id).Count(&dataCount).Error
		if err != nil {
			app.Error(c, -1, err, "查询服务类型失败")
			return
		}
		if dataCount > 0 {
			app.Error(c, -1, nil, "分组下存在服务分类，无法删除")
			return
		}
	}

	err = orm.Eloquent.Find(&oldData, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务分类失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = actions.AddAudit(c, tx, "业务", "服务分类", "删除", fmt.Sprintf("删除服务分类 \"%s\"", oldData.Name), oldData, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	err = tx.Delete(&business.ServiceClassify{}, id).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除分组列表失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 编辑服务分类
func EditServiceClassify(c *gin.Context) {
	var (
		err     error
		data    business.ServiceClassify
		oldData business.ServiceClassify
		id      string
	)

	id = c.Param("id")

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务分类失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = actions.AddAudit(c, tx, "业务", "服务分类", "编辑", fmt.Sprintf("编辑服务分类 \"%s\"", oldData.Name),
		map[string]interface{}{
			"identifies": oldData.Identifies,
			"name":       oldData.Name,
		},
		map[string]interface{}{
			"identifies": data.Identifies,
			"name":       data.Name,
		})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	err = tx.Model(&data).Where("id = ?", id).Updates(map[string]interface{}{
		"identifies": data.Identifies,
		"name":       data.Name,
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "编辑服务分类失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}
