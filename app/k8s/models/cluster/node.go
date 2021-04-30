package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

type Node struct {
	Id        int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name      string `gorm:"column:name;type:varchar(255)" json:"name"`
	HostId    int    `gorm:"column:host_id;type:int(11)" json:"host_id"`
	ClusterId int    `gorm:"column:cluster_id;type:int(11)" json:"cluster_id"`
	Role      string `gorm:"column:role;type:varchar(255)" json:"role"`
	Status    string `gorm:"column:status;type:varchar(255)" json:"status"`
	Message   string `gorm:"column:message;type:mediumtext" json:"message"`
	Dirty     int    `gorm:"column:dirty;type:int(11);default:0" json:"dirty"`
	models.BaseModel
}

func (Node) TableName() string {
	return "k8s_cluster_node"
}
