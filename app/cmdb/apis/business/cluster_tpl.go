package business

import (
	"fiy/app/cmdb/models/business"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools"
	"fiy/tools/app"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
  @Desc : 集群模板
*/

// 新建集群模板
func CreateClusterTpl(c *gin.Context) {
	var (
		err  error
		data struct {
			Name    string `json:"name"`
			SvcTpls []int  `json:"svc_tpls"`
		}
		svcReTpl []business.ServiceCluster
	)

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	userId := tools.GetUserId(c)

	CTData := business.ClusterTemplate{
		Name:     data.Name,
		Creator:  userId,
		Modifier: userId,
	}

	tx := orm.Eloquent.Begin()
	// 创建集群模板
	err = tx.Model(&business.ClusterTemplate{}).Create(&CTData).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	if len(data.SvcTpls) > 0 {
		for _, i := range data.SvcTpls {
			svcReTpl = append(svcReTpl, business.ServiceCluster{
				SvcTpl:     i,
				ClusterTpl: CTData.Id,
			})
		}
		err = tx.Create(&svcReTpl).Error
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, "新建集群模板失败")
			return
		}
	}

	// 添加操作审计
	err = actions.AddAudit(c, tx, "业务", "集群模版", "新建", fmt.Sprintf("新建集群模版 \"%s\"", data.Name), map[string]interface{}{}, data)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 集群模板列表
func ClusterTplList(c *gin.Context) {
	var (
		err    error
		result interface{}
		list   []struct {
			business.ClusterTemplate
			ModifyName string `json:"modify_name"`
		}
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	db := orm.Eloquent.Model(&business.ClusterTemplate{}).
		Joins("left join sys_user on sys_user.user_id = cmdb_business_cluster_tpl.modifier").
		Select("cmdb_business_cluster_tpl.*, sys_user.nick_name as modify_name")

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &list, SearchParams, "cmdb_business_cluster_tpl")
	if err != nil {
		app.Error(c, -1, err, "查询集群模板列表失败")
		return
	}

	app.OK(c, result, "")
}

// 集群模板关联的服务模板
func ClusterTplReSvcTpl(c *gin.Context) {
	var (
		err        error
		id         string
		svcTplList []business.ServiceTemplate
		ids        []int
	)

	id = c.Param("id")

	// 查询集群模板绑定的服务模板对应的ID
	err = orm.Eloquent.Model(&business.ServiceCluster{}).
		Where("cluster_tpl = ?", id).
		Pluck("svc_tpl", &ids).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务模板ID列表失败")
		return
	}

	// 查询ID列表对应的服务列表数据
	err = orm.Eloquent.Model(&business.ServiceTemplate{}).
		Where("id in (?)", ids).
		Find(&svcTplList).Error
	if err != nil {
		app.Error(c, -1, err, "查询服务列表失败")
		return
	}

	app.OK(c, svcTplList, "")
}

// 编辑集群模板
func EditClusterTpl(c *gin.Context) {
	var (
		err  error
		data struct {
			Name    string `json:"name"`
			SvcTpls []int  `json:"svc_tpls"`
		}
		id       string
		svcReTpl []business.ServiceCluster
	)

	id = c.Param("id")
	idInt, _ := strconv.Atoi(id)

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()
	// 查询变更前的数据
	var oldData struct {
		Name    string `json:"name"`
		SvcTpls []int  `json:"svc_tpls"`
	}
	err = tx.Model(&business.ClusterTemplate{}).Where("id = ?", id).Find(&oldData).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "查询变更前数据失败")
		return
	}

	err = tx.Model(&business.ServiceCluster{}).Where("cluster_tpl = ?", id).Pluck("svc_tpl", &oldData.SvcTpls).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "查询变更前数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c, tx, "业务", "集群模版", "编辑", fmt.Sprintf("编辑集群模版 \"%s\"", data.Name), oldData, data)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	// 编辑集群模板
	err = tx.Model(&business.ClusterTemplate{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":     data.Name,
		"modifier": tools.GetUserId(c),
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "编辑集群模板基础信息失败")
		return
	}

	// 编辑集群模板关联的服务模板数据
	err = tx.Where("cluster_tpl = ?", id).Delete(&business.ServiceCluster{}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除关联服务模板数据失败")
		return
	}
	if len(data.SvcTpls) > 0 {
		for _, i := range data.SvcTpls {
			svcReTpl = append(svcReTpl, business.ServiceCluster{
				SvcTpl:     i,
				ClusterTpl: idInt,
			})
		}
		err = tx.Create(&svcReTpl).Error
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, "新建集群模板失败")
			return
		}
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除集群模板
func DeleteClusterTpl(c *gin.Context) {
	var (
		err   error
		id    string
		count int64
	)

	id = c.Param("id")

	// 判断是否存在绑定的服务模板
	err = orm.Eloquent.Model(&business.ServiceCluster{}).Where("cluster_tpl = ?", id).Count(&count).Error
	if err != nil {
		app.Error(c, -1, err, "查询绑定的关联数据失败")
		return
	}
	if count > 0 {
		app.Error(c, -1, nil, "集群模板存在绑定的服务模板，无法删除")
		return
	}

	tx := orm.Eloquent.Begin()

	// 查询删除前的数据
	var oldData business.ClusterTemplate
	err = tx.Find(&oldData, id).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "查询删除前的数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c, tx, "业务", "集群模版", "删除", fmt.Sprintf("删除集群模版 \"%s\"", oldData.Name), oldData, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	// 删除集群模板
	err = tx.Delete(&business.ClusterTemplate{}, id).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, nil, "删除集群模板失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}
