package settings

import (
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

// 凭证
type Credential struct {
	Id       int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name; type:varchar(128);" json:"name" binding:"required"`
	UserName string `gorm:"column:username; type:varchar(128);" json:"username" binding:"required"`
	Type     int    `gorm:"column:type; type:int(11);" json:"type" binding:"required"` // 1 密码，2 密钥
	Content  string `gorm:"column:content; type:mediumtext;" json:"content" binding:"required"`
	models.BaseModel
}

func (Credential) TableName() string {
	return "k8s_setting_credential"
}
