package resource

import (
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 获取数据列表
func DataList(c *gin.Context) {
	var (
		err        error
		dataList   []*resource.Data
		result     interface{}
		modelId    string
		value      string
		identifies string
	)

	modelId = c.Param("id")

	db := orm.Eloquent.Model(&resource.Data{}).Where("info_id = ?", modelId)

	value = c.DefaultQuery("value", "")
	identifies = c.DefaultQuery("identifies", "")
	if identifies != "" && value != "" {
		db = db.Where(fmt.Sprintf("data->'$.%s' like '%%%s%%'", identifies, value))
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &dataList)
	if err != nil {
		app.Error(c, -1, err, "分页查询关联类型列表失败")
		return
	}

	app.OK(c, result, "")
}

// 新建数据
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

// 删除数据
func DeleteData(c *gin.Context) {
	var (
		err    error
		dataId string
	)

	dataId = c.Param("id")

	err = orm.Eloquent.Delete(&resource.Data{}, dataId).Error
	if err != nil {
		app.Error(c, -1, err, "删除资源数据失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑数据
func EditData(c *gin.Context) {
	var (
		err    error
		data   resource.Data
		dataId string
	)

	dataId = c.Param("id")

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&data).Where("id = ?", dataId).Updates(map[string]interface{}{
		"info_id": data.InfoId,
		"data":    data.Data,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "更新资源数据失败")
		return
	}

	app.OK(c, nil, "")
}
