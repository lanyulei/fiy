package version

import (
	"fiy/app/k8s/models/version"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 版本列表
func List(c *gin.Context) {
	var (
		err  error
		list []version.Manifest
	)

	err = orm.Eloquent.Model(&list).Order("id desc").Find(&list).Error
	if err != nil {
		app.Error(c, -1, err, "查询版本信息失败")
		return
	}

	app.OK(c, list, "")
}
