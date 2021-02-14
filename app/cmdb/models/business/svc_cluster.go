package business

/*
  @Author : lanyulei
  @Desc :
*/

type ServiceCluster struct {
	SvcTpl     int `gorm:"column:svc_tpl; type:int(11); uniqueIndex:idx_unique_svc_cluster;" json:"svc_tpl"`         // 服务模板
	ClusterTpl int `gorm:"column:cluster_tpl; type:int(11); uniqueIndex:idx_unique_svc_cluster;" json:"cluster_tpl"` // 集群模板
}

func (ServiceCluster) TableName() string {
	return "cmdb_business_svc_cluster"
}
