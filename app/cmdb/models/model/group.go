package model

import (
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

// 模型分组
type Group struct {
	Id         int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Identifies string `gorm:"column:identifies; type:varchar(128); unique;" json:"identifies" binding:"required"` // 分组唯一标识，只能是英文
	Name       string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"`             // 分组名称
	models.BaseModel
}

func (Group) TableName() string {
	return "cmdb_model_group"
}
