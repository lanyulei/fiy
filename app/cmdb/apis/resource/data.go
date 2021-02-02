package resource

import (
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 获取业务模型详情
func GetBizDetails(c *gin.Context) {
	app.OK(c, nil, "")
}

// 新建业务数据
func CreateData(c *gin.Context) {
	var (
		err  error
		data resource.Data
	)

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Create(&data).Error
	if err != nil {
		app.Error(c, -1, err, "新建资源数据失败")
		return
	}

	app.OK(c, nil, "")
}
