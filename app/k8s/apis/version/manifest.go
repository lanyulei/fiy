package version

import (
	"fiy/app/k8s/models/version"
	orm "fiy/common/global"
	"fiy/tools/app"
	"strings"

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

	content := make([]map[string][]version.Manifest, 0)
	tmpContent := make(map[string][]version.Manifest)
	tmpVersionList := make([]string, 0)
	for _, v := range list {
		vs := strings.Join(strings.Split(v.Version, ".")[0:2], ".")
		if _, ok := tmpContent[vs]; ok {
			tmpContent[vs] = append(tmpContent[vs], v)
		} else {
			tmpContent[vs] = []version.Manifest{v}
			tmpVersionList = append(tmpVersionList, vs)
		}
	}

	for _, v := range tmpVersionList {
		content = append(content, map[string][]version.Manifest{
			v: tmpContent[v],
		})
	}

	app.OK(c, content, "")
}
