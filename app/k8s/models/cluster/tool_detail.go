package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

type ToolDetail struct {
	Id           int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name         string `gorm:"column:name;type:varchar(255)" json:"name"`
	Version      string `gorm:"column:version;type:varchar(255)" json:"version"`
	ChartVersion string `gorm:"column:chart_version;type:varchar(255)" json:"chart_version"`
	Architecture string `gorm:"column:architecture;type:varchar(255)" json:"architecture"`
	Describe     string `gorm:"column:describe;type:varchar(255)" json:"describe"`
	Vars         string `gorm:"column:vars;type:mediumtext" json:"vars"`
	models.BaseModel
}

func (ToolDetail) TableName() string {
	return "k8s_cluster_tool_detail"
}
