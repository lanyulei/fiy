package settings

import (
	"fiy/app/k8s/models/settings"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建NTP服务
func SaveNTP(c *gin.Context) {
	var (
		err    error
		server settings.NTP
		count  int64
	)

	err = c.ShouldBind(&server)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&server).Count(&count).Error
	if err != nil {
		app.Error(c, -1, err, "查询NTP服务失败")
		return
	}
	if count == 0 {
		err = orm.Eloquent.Create(&server).Error
		if err != nil {
			app.Error(c, -1, err, "创建NTP服务失败")
			return
		}
	} else {
		err = orm.Eloquent.Model(&server).Update("server", server.Server).Error
		if err != nil {
			app.Error(c, -1, err, "更新NTP服务失败")
			return
		}
	}

	app.OK(c, "", "")
}

func GetNTP(c *gin.Context) {
	var (
		err    error
		server settings.NTP
	)

	err = orm.Eloquent.Find(&server).Error
	if err != nil {
		app.Error(c, -1, err, "查询NTP服务失败")
		return
	}

	app.OK(c, server, "")
}
