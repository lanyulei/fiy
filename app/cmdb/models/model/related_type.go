package model

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 模型关联类型
type RelatedType struct {
	Id             int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Identifies     string `gorm:"column:identifies; type:varchar(128); unique;" json:"identifies" binding:"required"`   // 分组唯一标识，只能是英文
	Name           string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"`               // 分组名称
	SourceDescribe string `gorm:"column:source_describe; type:varchar(1024)" json:"source_describe" binding:"required"` // 源->目标描述
	TargetDescribe string `gorm:"column:target_describe; type:varchar(1024)" json:"target_describe" binding:"required"` // 目标->源描述
	Direction      int    `gorm:"column:direction; type:int(11)" json:"direction" binding:"required"`                   // 是否有方向  1：有，源指向目标，2：双向，3：无方向
	models.BaseModel
}

func (RelatedType) TableName() string {
	return "cmdb_model_related_type"
}
