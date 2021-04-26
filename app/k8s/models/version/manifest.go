package version

import (
	"fiy/common/models"

	"gorm.io/datatypes"
)

/*
  @Author : lanyulei
*/

// k8s版本管理
type Manifest struct {
	Id       int            `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name     string         `gorm:"column:name; type:varchar(128);" json:"name" binding:"required"`
	Version  string         `gorm:"column:version; type:varchar(128);" json:"version"`
	Core     datatypes.JSON `gorm:"column:core; type:json;" json:"core"`
	Network  datatypes.JSON `gorm:"column:network; type:json;" json:"network"`
	Tool     datatypes.JSON `gorm:"column:tool; type:json;" json:"tool"`
	Storage  datatypes.JSON `gorm:"column:storage; type:json;" json:"storage"`
	Other    datatypes.JSON `gorm:"column:other; type:json;" json:"other"`
	IsActive bool           `gorm:"column:is_active; type:tinyint(1); default:1" json:"is_active"`
	models.BaseModel
}

func (Manifest) TableName() string {
	return "k8s_version_manifest"
}
