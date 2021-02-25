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

// 新建进程
func CreateProcess(c *gin.Context) {
	var (
		err     error
		process business.ServiceTemplateProcess
	)

	err = c.ShouldBind(&process)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Create(&process).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建进程失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"业务",
		"服务模版",
		"新建",
		fmt.Sprintf("新建进程 \"%s\"", process.Name),
		map[string]interface{}{},
		process)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 编辑进程
func EditProcess(c *gin.Context) {
	var (
		err     error
		id      string
		process business.ServiceTemplateProcess
		oldData business.ServiceTemplateProcess
	)

	id = c.Param("id")

	err = c.ShouldBind(&process)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询进程失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Model(&process).Where("id = ?", id).Save(&process).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "编辑进程失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"业务",
		"服务模版",
		"编辑",
		fmt.Sprintf("编辑进程 \"%s\"", process.Name),
		oldData,
		process)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除进程
func DeleteProcess(c *gin.Context) {
	var (
		err     error
		id      string
		oldData business.ServiceTemplateProcess
	)

	id = c.Param("id")

	err = orm.Eloquent.Find(&oldData, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询进程失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&business.ServiceTemplateProcess{}, id).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除进程失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"业务",
		"服务模版",
		"删除",
		fmt.Sprintf("删除进程 \"%s\"", oldData.Name),
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
