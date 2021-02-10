package business

import "fiy/common/models"

/*
  @Author : lanyulei
  @Desc : 服务模板
*/

type ServiceTemplate struct {
	Id          int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name        string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"` // 名称
	SvcClassify int    `gorm:"column:svc_classify; type:int(11);" json:"svc_classify"`                 // 服务分类
	Creator     int    `gorm:"column:creator; type:int(11);" json:"creator"`                           // 创建者
	Modifier    int    `gorm:"column:modifier; type:int(11);" json:"modifier"`                         // 修改者
	models.BaseModel
}

func (ServiceTemplate) TableName() string {
	return "cmdb_business_svc_tpl"
}
