package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 项目
type Info struct {
	Id     int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name   string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"` // 名称
	Spec   int    `gorm:"column:spec; type:int(11);" json:"spec"`                                 // 规格
	Status string `gorm:"column:status; type:varchar(45);" json:"status"`                         // 状态
	models.BaseModel
}

func (Info) TableName() string {
	return "k8s_cluster_info"
}
