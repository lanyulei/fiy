package cluster

import (
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

type Status struct {
	Id       int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Message  string `gorm:"column:message;type:mediumtext" json:"message"`
	Phase    string `gorm:"column:phase;type:varchar(255)" json:"phase"`
	PrePhase string `gorm:"column:pre_phase;type:varchar(255)" json:"pre_phase"`
	models.BaseModel
}

func (Status) TableName() string {
	return "k8s_cluster_status"
}
