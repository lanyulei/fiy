package model

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 模型关联类型
type RelatedType struct {
	Id             int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Identifies     string `gorm:"column:identifies; type:varchar(128)" json:"identifies"`            // 分组唯一标识，只能是英文
	Name           string `gorm:"column:name; type:varchar(128)" json:"name"`                        // 分组名称
	SourceDescribe string `gorm:"column:source_describe; type:varchar(1024)" json:"source_describe"` // 源->目标描述
	TargetDescribe string `gorm:"column:target_describe; type:varchar(1024)" json:"target_describe"` // 目标->源描述
	models.BaseModel
}

func (RelatedType) TableName() string {
	return "cmdb_model_related_type"
}
