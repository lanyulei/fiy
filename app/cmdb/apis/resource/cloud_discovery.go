package resource

import (
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建云资源同步任务
func CreateCloudDiscovery(c *gin.Context) {
	var (
		err       error
		userId    int
		discovery resource.CloudDiscovery
	)

	err = c.ShouldBind(&discovery)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	userId = tools.GetUserId(c)

	discovery.Creator = userId
	discovery.Modifier = userId

	err = orm.Eloquent.Create(&discovery).Error
	if err != nil {
		app.Error(c, -1, err, "创建云资源同步任务失败")
		return
	}

	app.OK(c, nil, "")
}

// 云资源同步任务列表
func CloudDiscoveryList(c *gin.Context) {
	var (
		err                error
		result             interface{}
		cloudDiscoveryList []*struct {
			resource.CloudDiscovery
			CreatorName      string `json:"creator_name"`
			ModifierName     string `json:"modifier_name"`
			CloudAccountName string `json:"cloud_account_name"`
			CloudAccountType string `json:"cloud_account_type"`
			ModelInfoName    string `json:"model_info_name"`
		}
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	db := orm.Eloquent.Model(&resource.CloudDiscovery{}).
		Joins("left join sys_user as suc on suc.user_id = cmdb_resource_cloud_discovery.creator").
		Joins("left join sys_user as sum on sum.user_id = cmdb_resource_cloud_discovery.modifier").
		Joins("left join cmdb_resource_cloud_account as ca on ca.id = cmdb_resource_cloud_discovery.cloud_account").
		Joins("left join cmdb_model_info as mi on mi.id = cmdb_resource_cloud_discovery.resource_type").
		Select("cmdb_resource_cloud_discovery.*, suc.nick_name as creator_name, sum.nick_name as modifier_name, ca.name as cloud_account_name, ca.type as cloud_account_type, mi.name as model_info_name")

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &cloudDiscoveryList, SearchParams, "cmdb_resource_cloud_discovery")
	if err != nil {
		app.Error(c, -1, err, "分页查询云资源同步失败")
		return
	}

	app.OK(c, result, "")
}

// 删除云资源同步任务
func DeleteCloudDiscovery(c *gin.Context) {
	var (
		err         error
		discoveryId string
	)

	discoveryId = c.Param("id")

	err = orm.Eloquent.Delete(&resource.CloudDiscovery{}, discoveryId).Error
	if err != nil {
		app.Error(c, -1, err, "删除云账号失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑云资源同步任务
func EditCloudDiscovery(c *gin.Context) {
	var (
		err         error
		discovery   resource.CloudDiscovery
		discoveryId string
	)

	discoveryId = c.Param("id")

	err = c.ShouldBind(&discovery)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&discovery).Where("id = ?", discoveryId).Updates(map[string]interface{}{
		"name":          discovery.Name,
		"resource_type": discovery.ResourceType,
		"cloud_account": discovery.CloudAccount,
		"field_map":     discovery.FieldMap,
		"status":        discovery.Status,
		"modifier":      tools.GetUserId(c),
		"remarks":       discovery.Remarks,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "编辑云资源同步任务失败")
		return
	}

	app.OK(c, nil, "")
}
