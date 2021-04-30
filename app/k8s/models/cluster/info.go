package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 项目
type Info struct {
	Id       int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"` // 名称
	Source   string `gorm:"column:source;type:varchar(255)" json:"source"`
	SpecId   int    `gorm:"column:spec_id;type:int(11)" json:"spec_id"`
	SecretId int    `gorm:"column:secret_id;type:int(11)" json:"secret_id"`
	StatusId int    `gorm:"column:status_id;type:int(11)" json:"status_id"`
	PlanId   int    `gorm:"column:plan_id;type:int(11)" json:"plan_id"`
	LogId    int    `gorm:"column:log_id;type:int(11)" json:"log_id"`
	Dirty    int    `gorm:"column:dirty;type:int(11);default:0" json:"dirty"`
	models.BaseModel
}

func (Info) TableName() string {
	return "k8s_cluster_info"
}
