package analysis

import (
	"fiy/app/cmdb/models/business"
	"fiy/app/cmdb/models/model"
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 运营分析
func Operation(c *gin.Context) {
	var (
		err             error
		bizCount        int64
		modelCount      int64
		clusterTplCount int64
		svcTplCount     int64
	)

	// 业务总数
	err = orm.Eloquent.Model(&resource.Data{}).
		Joins("left join cmdb_model_info on cmdb_model_info.id = cmdb_resource_data.info_id").
		Where("cmdb_model_info.identifies = 'built_in_biz'").
		Count(&bizCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询业务总数失败")
		return
	}

	// 模型数量
	err = orm.Eloquent.Model(&model.Info{}).Count(&modelCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型总数失败")
		return
	}

	// 集群模版总数
	err = orm.Eloquent.Model(&business.ClusterTemplate{}).Count(&clusterTplCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询集群模板总数失败")
		return
	}

	// 服务模版总数
	err = orm.Eloquent.Model(&business.ServiceTemplate{}).Count(&svcTplCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询集群模板总数失败")
		return
	}

	result := map[string]interface{}{
		"biz_count":         bizCount,
		"model_count":       modelCount,
		"cluster_tpl_count": clusterTplCount,
		"svc_tpl_count":     svcTplCount,
	}

	// 所有模型的资源统计

	app.OK(c, result, "")
}
