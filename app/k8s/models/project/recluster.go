package project

import "fiy/common/models"

/*
  @Author : lanyulei
*/

type ReCluster struct {
	Id        int `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	ProjectID int `gorm:"column:project_id;type:int(11)" json:"project_id"`
	ClusterID int `gorm:"column:cluster_id;type:int(11)" json:"cluster_id"`
	models.BaseModel
}

func (ReCluster) TableName() string {
	return "k8s_project_cluster"
}
