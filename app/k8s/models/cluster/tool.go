package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

type Tool struct {
	Id            int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name          string `gorm:"column:name;type:varchar(255)" json:"name"`
	ClusterId     string `gorm:"column:cluster_id;type:varchar(255)" json:"cluster_id"`
	Version       string `gorm:"column:version;type:varchar(255)" json:"version"`
	Describe      string `gorm:"column:describe;type:varchar(255)" json:"describe"`
	Status        string `gorm:"column:status;type:varchar(255)" json:"status"`
	Message       string `gorm:"column:message;type:mediumtext" json:"message"`
	Logo          string `gorm:"column:logo;type:varchar(255)" json:"logo"`
	Vars          string `gorm:"column:vars;type:mediumtext" json:"vars"`
	Url           string `gorm:"column:url;type:varchar(255)" json:"url"`
	Frame         int    `gorm:"column:frame;type:int(11);default:0" json:"frame"`
	Architecture  string `gorm:"column:architecture;type:varchar(255);default:all" json:"architecture"`
	HigherVersion string `gorm:"column:higher_version;type:varchar(255)" json:"higher_version"`
	models.BaseModel
}

func (Tool) TableName() string {
	return "k8s_cluster_tool"
}
