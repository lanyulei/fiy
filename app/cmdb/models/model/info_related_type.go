package model

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 模型和关联类型的关联表
type InfoRelatedType struct {
	Id            int `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	InfoId        int `gorm:"column:info_id; type:int(11); uniqueIndex:idx_related_info;" json:"info_id"`                 // 模型ID
	RelatedTypeId int `gorm:"column:related_type_id; type:int(11); uniqueIndex:idx_related_info;" json:"related_type_id"` // 关联模型ID
	models.BaseModel
}

func (InfoRelatedType) TableName() string {
	return "cmdb_model_info_related_type"
}
