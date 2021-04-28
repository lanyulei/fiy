package settings

import (
	"fiy/app/k8s/models/settings"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建仓库信息
func CreateRegistry(c *gin.Context) {
	var (
		err      error
		registry settings.Registry
	)

	err = c.ShouldBind(&registry)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Create(&registry).Error
	if err != nil {
		app.Error(c, -1, err, "新建仓库信息失败")
		return
	}

	app.OK(c, "", "")
}

// 编辑仓库信息
func EditRegistry(c *gin.Context) {
	var (
		err      error
		registry settings.Registry
		id       string
	)

	id = c.Param("id")

	err = c.ShouldBind(&registry)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&registry).Where("id = ?", id).Updates(map[string]interface{}{
		"ip":           registry.IP,
		"protocol":     registry.Protocol,
		"architecture": registry.Architecture,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "编辑仓库数据失败")
		return
	}

	app.OK(c, "", "")
}

// 仓库信息列表
func RegistryList(c *gin.Context) {
	var (
		err    error
		list   []settings.Registry
		result interface{}
		ip     string
	)
	db := orm.Eloquent.Model(&settings.Registry{})
	ip = c.DefaultQuery("ip", "")
	if ip != "" {
		db = db.Where("ip like ?", "%"+ip+"%")
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &list)
	if err != nil {
		app.Error(c, -1, err, "查询仓库信息列表失败")
		return
	}

	app.OK(c, result, "")
}

// 删除仓库信息
func DeleteRegistry(c *gin.Context) {
	var (
		err error
		id  string
	)

	id = c.Param("id")

	err = orm.Eloquent.Delete(&settings.Registry{}, id).Error
	if err != nil {
		app.Error(c, -1, err, "删除仓库信息失败")
		return
	}

	app.OK(c, "", "")
}
