package business

import "fiy/common/models"

/*
  @Author : lanyulei
  @Desc : 集群模板
*/

type ClusterTemplate struct {
	Id       int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"` // 名称
	Creator  int    `gorm:"column:creator; type:int(11);" json:"creator"`                           // 创建者
	Modifier int    `gorm:"column:modifier; type:int(11);" json:"modifier"`                         // 修改者
	models.BaseModel
}

func (ClusterTemplate) TableName() string {
	return "cmdb_business_cluster_tpl"
}
