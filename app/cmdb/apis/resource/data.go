package resource

import (
	"fiy/app/cmdb/models/resource"
	"fiy/common/actions"
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

	tx := orm.Eloquent.Begin()

	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"新建",
		fmt.Sprintf("新建资源数据 <%d>", data.Id),
		map[string]interface{}{},
		data)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除数据
func DeleteData(c *gin.Context) {
	var (
		err    error
		dataId string
		data   resource.Data
	)

	dataId = c.Param("id")

	err = orm.Eloquent.Find(&data, dataId).Error
	if err != nil {
		app.Error(c, -1, err, "查询资源数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&resource.Data{}, dataId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"删除",
		fmt.Sprintf("删除资源数据 \"%s\"", dataId),
		data,
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 编辑数据
func EditData(c *gin.Context) {
	var (
		err     error
		data    resource.Data
		oldData resource.Data
		dataId  string
	)

	dataId = c.Param("id")

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, dataId).Error
	if err != nil {
		app.Error(c, -1, err, "查询资源数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Model(&data).Where("id = ?", dataId).Updates(map[string]interface{}{
		"info_id": data.InfoId,
		"data":    data.Data,
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"编辑",
		fmt.Sprintf("编辑资源数据 \"%s\"", dataId),
		map[string]interface{}{
			"info_id": oldData.InfoId,
			"data":    oldData.Data,
		},
		map[string]interface{}{
			"info_id": data.InfoId,
			"data":    data.Data,
		})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 获取数据详情
func GetDataDetails(c *gin.Context) {
	var (
		err     error
		details resource.Data
		dataId  string
	)

	dataId = c.Param("id")

	err = orm.Eloquent.Model(&details).Where("id = ?", dataId).Find(&details).Error
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	app.OK(c, details, "")
}
