package business

import (
	"fiy/app/cmdb/models/business"
	orm "fiy/common/global"
	"fiy/tools/app"

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

	err = orm.Eloquent.Create(&process).Error
	if err != nil {
		app.Error(c, -1, err, "新建进程失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑进程
func EditProcess(c *gin.Context) {
	var (
		err     error
		id      string
		process business.ServiceTemplateProcess
	)

	id = c.Param("id")

	err = c.ShouldBind(&process)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&process).Where("id = ?", id).Save(&process).Error
	if err != nil {
		app.Error(c, -1, err, "编辑进程失败")
		return
	}

	app.OK(c, nil, "")
}

// 删除进程
func DeleteProcess(c *gin.Context) {
	var (
		err error
		id  string
	)

	id = c.Param("id")

	err = orm.Eloquent.Delete(&business.ServiceTemplateProcess{}, id).Error
	if err != nil {
		app.Error(c, -1, err, "删除进程失败")
		return
	}

	app.OK(c, nil, "")
}
