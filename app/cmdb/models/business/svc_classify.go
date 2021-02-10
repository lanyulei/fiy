package business

import "fiy/common/models"

/*
  @Author : lanyulei
  @Desc : 服务分类
*/

type ServiceClassify struct {
	Id         int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Identifies string `gorm:"column:identifies; type:varchar(128); unique;" json:"identifies" binding:"required"` // 标识
	Name       string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"`             // 名称
	Level      int    `gorm:"column:level; type:int(11);" json:"level"`                                           // 级别 1 分组，2 服务
	Parent     int    `gorm:"column:parent; type:int(11);" json:"parent"`                                         // 父级
	models.BaseModel
}

func (ServiceClassify) TableName() string {
	return "cmdb_business_svc_classify"
}
