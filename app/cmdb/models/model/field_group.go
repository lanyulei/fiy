package model

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 字段分组
type FieldGroup struct {
	Id       int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name; type:varchar(50); uniqueIndex:idx_name_info_id;" json:"name" binding:"required"` // 分组名称
	Sequence int    `gorm:"column:sequence; type:int(11)" json:"sequence"`                                               // 分组展示的顺序
	IsFold   bool   `gorm:"column:is_fold; type:tinyint(1)" json:"is_fold"`                                              // 是否折叠
	InfoId   int    `gorm:"column:info_id; type:int(11); uniqueIndex:idx_name_info_id;" json:"info_id"`                  // 对应的模型ID
	models.BaseModel
}

func (FieldGroup) TableName() string {
	return "cmdb_model_field_group"
}
