package resource

/*
  @Author : lanyulei
  @Desc : 数据关联表
*/

type DataRelated struct {
	Source int `gorm:"column:source; index; uniqueIndex:idx_source_target" json:"source"` // 源ID
	Target int `gorm:"column:target; index; uniqueIndex:idx_source_target" json:"target"` // 目标ID
}

func (DataRelated) TableName() string {
	return "cmdb_resource_data_related"
}
