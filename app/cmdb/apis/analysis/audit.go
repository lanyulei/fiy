package analysis

import (
	"fiy/app/cmdb/models/analysis"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
  @Desc : 操作审计
*/

// 审计列表
func AuditList(c *gin.Context) {
	var (
		err    error
		result interface{}
		list   []analysis.Audit
	)

	db := orm.Eloquent.Model(&analysis.Audit{})

	field := c.DefaultQuery("field", "")
	value := c.DefaultQuery("value", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")
	if field != "" && value != "" {
		db = db.Where(fmt.Sprintf("%s like '%s'", field, "%"+value+"%"))
	}

	if startTime != "" && endTime != "" {
		db = db.Where("created_at between ? and ?", startTime, endTime)
	}

	db = db.Select("cmdb_analysis_audit.id, cmdb_analysis_audit.classify, cmdb_analysis_audit.features, cmdb_analysis_audit.action, cmdb_analysis_audit.describe, cmdb_analysis_audit.created_at, cmdb_analysis_audit.username")
	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &list)
	if err != nil {
		app.Error(c, -1, err, "查询操作历史失败")
		return
	}

	app.OK(c, result, "")
}

// 详情
func AuditDetails(c *gin.Context) {
	var (
		err  error
		id   string
		info analysis.Audit
	)

	id = c.Param("id")

	err = orm.Eloquent.Find(&info, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询审计数据详情失败")
		return
	}

	app.OK(c, info, "")
}
