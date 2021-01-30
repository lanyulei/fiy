package model

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 模型和关联类型的关联表
type InfoRelatedType struct {
	Id            int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Source        int    `gorm:"column:source; type:int(11); uniqueIndex:idx_unique;" json:"source"`                   // 源模型ID
	Target        int    `gorm:"column:target; type:int(11); uniqueIndex:idx_unique;" json:"target"`                   // 目标模型ID
	RelatedTypeId int    `gorm:"column:related_type_id; type:int(11); uniqueIndex:idx_unique;" json:"related_type_id"` // 关联类型ID
	Constraint    int    `gorm:"column:constraint; type:int(11); uniqueIndex:idx_unique;" json:"constraint"`           // 源-目标约束
	Remark        string `gorm:"column:remark; type:varchar(1024);" json:"remark"`                                     // 关联描述
	models.BaseModel
}

func (InfoRelatedType) TableName() string {
	return "cmdb_model_info_related_type"
}
