package settings

import (
	"fiy/app/k8s/models/settings"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 新建凭据
func CreateCredential(c *gin.Context) {
	var (
		err        error
		credential settings.Credential
	)

	err = c.ShouldBind(&credential)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Create(&credential).Error
	if err != nil {
		app.Error(c, -1, err, "新建凭据失败")
		return
	}

	app.OK(c, "", "")
}

// 编辑凭据
func EditCredential(c *gin.Context) {
	var (
		err        error
		credential settings.Credential
		id         string
	)

	id = c.Param("id")

	err = c.ShouldBind(&credential)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Model(&credential).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":     credential.Name,
			"username": credential.Username,
			"type":     credential.Type,
			"content":  credential.Content,
		}).Error
	if err != nil {
		app.Error(c, -1, err, "编辑凭据失败")
		return
	}

	app.OK(c, "", "")
}

// 查询凭据列表
func CredentialList(c *gin.Context) {
	var (
		err      error
		list     []settings.Credential
		result   interface{}
		username string
	)
	db := orm.Eloquent.Model(&settings.Credential{})
	username = c.DefaultQuery("ip", "")
	if username != "" {
		db = db.Where("username like ?", "%"+username+"%")
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &list)
	if err != nil {
		app.Error(c, -1, err, "查询凭据列表失败")
		return
	}

	app.OK(c, result, "")
}

// 删除凭据
func DeleteCredential(c *gin.Context) {
	var (
		err error
		id  string
	)

	id = c.Param("id")

	err = orm.Eloquent.Delete(&settings.Credential{}, id).Error
	if err != nil {
		app.Error(c, -1, err, "删除凭据失败")
		return
	}

	app.OK(c, "", "")
}
