package analysis

import (
	"fiy/common/models"

	"gorm.io/datatypes"
)

/*
  @Author : lanyulei
  @Desc : 操作审计
*/

type Audit struct {
	Id       int            `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Classify string         `gorm:"column:classify; type:varchar(45);" json:"classify"`   // 分类
	Features string         `gorm:"column:features; type:varchar(45);" json:"features"`   // 功能模块
	Action   string         `gorm:"column:action; type:varchar(45);" json:"action"`       // 动作
	Describe string         `gorm:"column:describe; type:varchar(1024);" json:"describe"` // 描述
	Username string         `gorm:"column:username; type:varchar(128);" json:"username"`  // 操作账号
	OldData  datatypes.JSON `gorm:"column:old_data; type:json;" json:"old_data"`          // 变更前数据
	NewData  datatypes.JSON `gorm:"column:new_data; type:json;" json:"new_data"`          // 变更后数据
	models.BaseModel
}

func (Audit) TableName() string {
	return "cmdb_analysis_audit"
}
