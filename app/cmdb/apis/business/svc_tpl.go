package business

import (
	"fiy/app/cmdb/models/business"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/common/models"
	"fiy/common/pagination"
	"fiy/tools"
	"fiy/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 服务模板列表
func ServiceTemplateList(c *gin.Context) {
	var (
		err    error
		result interface{}
		list   []struct {
			business.ServiceTemplate
			SvcClassifyName string `json:"svc_classify_name"`
			ModifyName      string `json:"modify_name"`
		}
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	db := orm.Eloquent.Model(&business.ServiceTemplate{}).
		Joins("left join cmdb_business_svc_classify as sc on sc.id = cmdb_business_svc_tpl.svc_classify").
		Joins("left join sys_user on sys_user.user_id = cmdb_business_svc_tpl.modifier").
		Select("cmdb_business_svc_tpl.*, sc.name as svc_classify_name, sys_user.nick_name as modify_name")

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &list, SearchParams, "cmdb_business_svc_tpl")
	if err != nil {
		app.Error(c, -1, err, "分页查询云账号失败")
		return
	}

	app.OK(c, result, "")
}

// 新建服务模板
func CreateServiceTemplate(c *gin.Context) {
	var (
		err    error
		params struct {
			business.ServiceTemplate
			ProcessList []*business.ServiceTemplateProcess `json:"process_list"`
		}
	)

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 新建服务模板
	currentUser := tools.GetUserId(c)
	svcData := business.ServiceTemplate{
		Name:        params.Name,
		SvcClassify: params.SvcClassify,
		Creator:     currentUser,
		Modifier:    currentUser,
		BaseModel:   models.BaseModel{},
	}
	err = tx.Create(&svcData).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建服务模板失败")
		return
	}

	// 新建服务模板的进程
	if len(params.ProcessList) > 0 {
		for _, p := range params.ProcessList {
			p.SvcTpl = svcData.Id
		}
		err = tx.Create(&params.ProcessList).Error
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, "新建服务模板失败")
			return
		}
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"业务",
		"服务模版",
		"新建",
		fmt.Sprintf("新建服务模版 \"%s\"", params.Name),
		map[string]interface{}{},
		params)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 服务模板详情
func ServiceTemplateDetails(c *gin.Context) {
	var (
		err  error
		id   string
		info struct {
			business.ServiceTemplate
			ProcessList []*business.ServiceTemplateProcess `json:"process_list"`
		}
	)

	id = c.Param("id")

	// 查询服务模板
	err = orm.Eloquent.Model(&business.ServiceTemplate{}).Where("id = ?", id).Find(&info).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务模板失败")
		return
	}

	// 查询服务模板进程
	err = orm.Eloquent.Model(&business.ServiceTemplateProcess{}).
		Where("svc_tpl = ?", id).
		Find(&info.ProcessList).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务模板进程列表失败")
		return
	}

	app.OK(c, info, "")
}

// 删除服务模板
func DeleteServiceTemplate(c *gin.Context) {
	var (
		err          error
		id           string
		processCount int64
		oldData      business.ServiceTemplate
	)

	id = c.Param("id")

	// 有进程数据，则无法删除服务模板
	err = orm.Eloquent.Model(&business.ServiceTemplateProcess{}).Where("svc_tpl = ?", id).Count(&processCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务模板进程数量失败")
		return
	}
	if processCount > 0 {
		app.Error(c, -1, err, "无法删除当前服务模板，存在绑定的进程")
		return
	}

	err = orm.Eloquent.Find(&oldData, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务模板失败")
		return
	}

	tx := orm.Eloquent.Begin()

	// 删除模板
	err = tx.Delete(&business.ServiceTemplate{}, id).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除服务模板失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"业务",
		"服务模版",
		"删除",
		fmt.Sprintf("删除服务模版 \"%s\"", oldData.Name),
		oldData,
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}
