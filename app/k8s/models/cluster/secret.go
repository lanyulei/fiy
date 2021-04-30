package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

type Secret struct {
	Id              int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	KubeadmToken    string `gorm:"column:kubeadm_token;type:mediumtext" json:"kubeadm_token"`
	KubernetesToken string `gorm:"column:kubernetes_token;type:mediumtext" json:"kubernetes_token"`
	models.BaseModel
}

func (Secret) TableName() string {
	return "k8s_cluster_secret"
}
