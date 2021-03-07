package business

/*
  @Author : lanyulei
  @Desc : 集群模板
*/

type TemplateRelated struct {
	Id          int `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	TplClassify int `gorm:"column:tpl_classify; index; uniqueIndex:idx_tpl_data" json:"tpl_classify"` // 模板类型 1 集群模板， 2 服务模板
	TplId       int `gorm:"column:tpl_id; index; uniqueIndex:idx_tpl_data" json:"tpl_id"`
	DataID      int `gorm:"column:data_id; index; uniqueIndex:idx_tpl_data" json:"data_id"`
	InfoId      int `gorm:"column:info_id; index;" json:"info_id"`
}

func (TemplateRelated) TableName() string {
	return "cmdb_business_tpl_related"
}
