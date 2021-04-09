package resource

import "fiy/common/models"

/*
  @Author : lanyulei
  @Desc : 数据关联表
*/

type DataRelated struct {
	Id           int `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Source       int `gorm:"column:source; index; uniqueIndex:idx_source_target" json:"source"` // 源ID
	Target       int `gorm:"column:target; index; uniqueIndex:idx_source_target" json:"target"` // 目标ID
	SourceInfoId int `gorm:"column:source_info_id; index;" json:"source_info_id"`               // 源模型ID
	TargetInfoId int `gorm:"column:target_info_id; index;" json:"target_info_id"`               // 目标模型ID
	models.BaseModel
}

func (DataRelated) TableName() string {
	return "cmdb_resource_data_related"
}
