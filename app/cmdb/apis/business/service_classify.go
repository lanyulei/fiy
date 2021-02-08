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
