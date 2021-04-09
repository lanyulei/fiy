package resource

import (
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
  @Desc : 数据关联表
*/

// 资产关联绑定
func DataRelated(c *gin.Context) {
	var (
		err          error
		relatedValue struct {
			Asset  map[string]interface{} `json:"asset"`
			Models map[string]interface{} `json:"models"`
		}
		relatedList []resource.DataRelated
	)

	err = c.ShouldBind(&relatedValue)
	if err != nil {
		app.Error(c, -1, err, "绑定参数失败")
		return
	}

	relatedList = make([]resource.DataRelated, 0)
	for _, a := range relatedValue.Asset["list"].([]interface{}) {
		for _, m := range relatedValue.Models["list"].([]interface{}) {
			relatedList = append(relatedList, resource.DataRelated{
				Source:       int(m.(float64)),
				Target:       int(a.(float64)),
				SourceInfoId: int(relatedValue.Models["model"].(float64)),
				TargetInfoId: int(relatedValue.Asset["model"].(float64)),
			})
		}
	}

	err = orm.Eloquent.Create(&relatedList).Error
	if err != nil {
		app.Error(c, -1, err, "创建数据关联失败")
		return
	}

	app.OK(c, "", "")
}

// 删除资产关联
func DeleteDataRelated(c *gin.Context) {
	var (
		err       error
		relatedID string
	)

	relatedID = c.Param("id")

	err = orm.Eloquent.Delete(&resource.DataRelated{}, relatedID).Error
	if err != nil {
		app.Error(c, -1, err, "删除数据关联失败")
		return
	}

	app.OK(c, "", "")
}
