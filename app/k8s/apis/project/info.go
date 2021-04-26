package project

import (
	"fiy/app/k8s/models/project"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 项目列表
func List(c *gin.Context) {
	var (
		err    error
		result interface{}
		list   []project.Info
	)

	db := orm.Eloquent.Model(&project.Info{})

	name := c.DefaultQuery("name", "")
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &list)
	if err != nil {
		app.Error(c, -1, err, "查询项目列表失败")
		return
	}

	app.OK(c, result, "")
}

// 创建项目
func Create(c *gin.Context) {
	var (
		err  error
		info project.Info
	)

	err = c.ShouldBind(&info)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Create(&info).Error
	if err != nil {
		app.Error(c, -1, err, "创建项目失败")
		return
	}

	app.OK(c, "", "成功")
}

// 编辑项目
func Edit(c *gin.Context) {
	var (
		err  error
		info project.Info
		id   string
	)

	id = c.Param("id")

	err = c.ShouldBind(&info)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&info).Where("id = ?", id).Updates(map[string]interface{}{
		"name":    info.Name,
		"remarks": info.Remarks,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "更新项目失败")
		return
	}

	app.OK(c, "", "成功")
}

// 删除项目
func Delete(c *gin.Context) {
	var (
		err error
		id  string
	)

	id = c.Param("id")

	// todo 需确认是否有集群，若存在集群，则该项目不允许删除

	err = orm.Eloquent.Delete(&project.Info{}, id).Error
	if err != nil {
		app.Error(c, -1, err, "删除项目失败")
		return
	}

	app.OK(c, "", "成功")
}
