package settings

import (
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

// 仓库
type Registry struct {
	Id           int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	IP           string `gorm:"column:ip; type:varchar(128);" json:"ip" binding:"required"`
	Protocol     string `gorm:"column:protocol; type:varchar(45);" json:"protocol"`
	Architecture string `gorm:"column:architecture; type:varchar(45);" json:"architecture"`
	models.BaseModel
}

func (Registry) TableName() string {
	return "k8s_setting_registry"
}
