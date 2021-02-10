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

// 创建服务分类
func CreateServiceClassify(c *gin.Context) {
	var (
		err             error
		serviceClassify business.ServiceClassify
	)

	err = c.ShouldBind(&serviceClassify)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 写入数据库
	err = orm.Eloquent.Create(&serviceClassify).Error
	if err != nil {
		app.Error(c, -1, err, "创建服务分类失败")
		return
	}

	app.OK(c, nil, "")
}

// 服务分类列表
func ServiceClassifyList(c *gin.Context) {
	var (
		err  error
		list []*struct {
			business.ServiceClassify
			Services []business.ServiceClassify
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

	err = orm.Eloquent.Delete(&business.ServiceClassify{}, id).Error
	if err != nil {
		app.Error(c, -1, err, "删除分组列表失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑服务分类
func EditServiceClassify(c *gin.Context) {
	var (
		err  error
		data business.ServiceClassify
		id   string
	)

	id = c.Param("id")

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&data).Where("id = ?", id).Updates(map[string]interface{}{
		"identifies": data.Identifies,
		"name":       data.Name,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "编辑服务分类失败")
		return
	}

	app.OK(c, nil, "")
}
