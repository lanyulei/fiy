package resource

/*
  @Author : lanyulei
  @Desc : 数据关联表
*/

type DataRelated struct {
	Source       int `gorm:"column:source; index; uniqueIndex:idx_source_target" json:"source"` // 源ID
	Target       int `gorm:"column:target; index; uniqueIndex:idx_source_target" json:"target"` // 目标ID
	SourceInfoId int `gorm:"column:source_info_id; index;" json:"source_info_id"`               // 模型ID
	TargetInfoId int `gorm:"column:target_info_id; index;" json:"target_info_id"`               // 模型ID
}

func (DataRelated) TableName() string {
	return "cmdb_resource_data_related"
}
