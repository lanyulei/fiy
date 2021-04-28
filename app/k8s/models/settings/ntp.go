package settings

import (
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

// 仓库
type NTP struct {
	Id     int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Server string `gorm:"column:server; type:varchar(128);" json:"server" binding:"required"`
	models.BaseModel
}

func (NTP) TableName() string {
	return "k8s_setting_ntp"
}
