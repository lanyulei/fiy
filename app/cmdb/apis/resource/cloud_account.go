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

// 新建云账户
func CreateCloudAccount(c *gin.Context) {
	var (
		err     error
		userId  int
		account resource.CloudAccount
	)

	err = c.ShouldBind(&account)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	userId = tools.GetUserId(c)

	account.Creator = userId
	account.Modifier = userId

	err = orm.Eloquent.Create(&account).Error
	if err != nil {
		app.Error(c, -1, err, "创建云账户失败")
		return
	}

	app.OK(c, nil, "")
}

// 云账户列表
func CloudAccountList(c *gin.Context) {
	var (
		err              error
		result           interface{}
		cloudAccountList []*struct {
			resource.CloudAccount
			CreatorName  string `json:"creator_name"`
			ModifierName string `json:"modifier_name"`
		}
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	db := orm.Eloquent.Model(&resource.CloudAccount{}).
		Joins("left join sys_user as suc on suc.user_id = cmdb_resource_cloud_account.creator").
		Joins("left join sys_user as sum on sum.user_id = cmdb_resource_cloud_account.modifier").
		Select("cmdb_resource_cloud_account.*, suc.nick_name as creator_name, sum.nick_name as modifier_name")

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &cloudAccountList, SearchParams, "cmdb_resource_cloud_account")
	if err != nil {
		app.Error(c, -1, err, "分页查询云账号失败")
		return
	}

	app.OK(c, result, "")
}

// 删除云账户
func DeleteCloudAccount(c *gin.Context) {
	var (
		err       error
		accountId string
	)

	accountId = c.Param("id")

	err = orm.Eloquent.Delete(&resource.CloudAccount{}, accountId).Error
	if err != nil {
		app.Error(c, -1, err, "删除云账号失败")
		return
	}

	app.OK(c, nil, "")
}

// 编辑云账户
func EditCloudAccount(c *gin.Context) {
	var (
		err       error
		account   resource.CloudAccount
		accountId string
	)

	accountId = c.Param("id")

	err = c.ShouldBind(&account)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&account).Where("id = ?", accountId).Updates(map[string]interface{}{
		"name":     account.Name,
		"type":     account.Type,
		"status":   account.Status,
		"secret":   account.Secret,
		"key":      account.Key,
		"modifier": tools.GetUserId(c),
		"remarks":  account.Remarks,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "编辑云账户失败")
		return
	}

	app.OK(c, nil, "")
}
