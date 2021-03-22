package resource

import (
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

// 云账户管理
type CloudAccount struct {
	Id       int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"`                     // 账户名称
	Type     string `gorm:"column:type; type:varchar(45); uniqueIndex:idx_unique;" json:"type" binding:"required"`      // 账户类型
	Status   bool   `gorm:"column:status; type:tinyint(1); default:1" json:"status"`                                    // 账号状态
	Secret   string `gorm:"column:secret; type:varchar(128); uniqueIndex:idx_unique;" json:"secret" binding:"required"` // accessSecret
	Key      string `gorm:"column:key; type:varchar(128); uniqueIndex:idx_unique;" json:"key" binding:"required"`       // accessKeyId
	Creator  int    `gorm:"column:creator; type:int(11);" json:"creator"`                                               // 创建者
	Modifier int    `gorm:"column:modifier; type:int(11);" json:"modifier"`                                             // 修改者
	Remarks  string `gorm:"column:remarks; type:varchar(1024);" json:"remarks"`                                         // 备注
	models.BaseModel
}

func (CloudAccount) TableName() string {
	return "cmdb_resource_cloud_account"
}
