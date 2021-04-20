package project

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 项目
type Info struct {
	Id      int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name    string `gorm:"column:name; type:varchar(128); unique;" json:"name" binding:"required"` // 名称
	Remarks string `gorm:"column:remarks; type:varchar(1024);" json:"remarks"`                     // 描述
	models.BaseModel
}

func (Info) TableName() string {
	return "k8s_project_info"
}
