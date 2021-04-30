package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

type Istio struct {
	Id        int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name      string `gorm:"column:name;type:varchar(255)" json:"name"`
	ClusterId string `gorm:"column:cluster_id;type:varchar(255)" json:"cluster_id"`
	Version   string `gorm:"column:version;type:varchar(255)" json:"version"`
	Describe  string `gorm:"column:describe;type:varchar(255)" json:"describe"`
	Status    string `gorm:"column:status;type:varchar(255)" json:"status"`
	Message   string `gorm:"column:message;type:mediumtext" json:"message"`
	Vars      string `gorm:"column:vars;type:mediumtext" json:"vars"`
	models.BaseModel
}

func (Istio) TableName() string {
	return "k8s_cluster_istio"
}
