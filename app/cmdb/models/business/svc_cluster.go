package business

import "fiy/common/models"

/*
  @Author : lanyulei
  @Desc :
*/

type ServiceCluster struct {
	Id         int `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	SvcTpl     int `gorm:"column:svc_tpl; type:int(11);" json:"svc_tpl"`         // 服务模板
	ClusterTpl int `gorm:"column:cluster_tpl; type:int(11);" json:"cluster_tpl"` // 集群模板
	models.BaseModel
}

func (ServiceCluster) TableName() string {
	return "cmdb_business_svc_cluster"
}
