package actions

import (
	"encoding/json"
	"fiy/app/cmdb/models/analysis"
	orm "fiy/common/global"
	"fiy/tools"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 添加审计
func AddAudit(c *gin.Context, classify, features, action, describe string, oldData, newData interface{}) (err error) {
	/*
		classify 分组
		features 功能模块
		action 动作
		describe 描述
		username 操作账号
		oldData 变更前数据
		newData 变更后数据
	*/

	var (
		oldDatabyte []byte
		newDatabyte []byte
	)

	oldDatabyte, err = json.Marshal(oldData)
	if err != nil {
		return
	}

	newDatabyte, err = json.Marshal(newData)
	if err != nil {
		return
	}

	err = orm.Eloquent.Create(&analysis.Audit{
		Classify: classify,
		Features: features,
		Action:   action,
		Describe: describe,
		Username: fmt.Sprintf("%s (%s)", tools.GetNickName(c), tools.GetUserName(c)),
		OldData:  oldDatabyte,
		NewData:  newDatabyte,
	}).Error

	return
}
