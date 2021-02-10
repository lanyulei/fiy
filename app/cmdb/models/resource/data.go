package resource

import (
	"fiy/common/models"

	"gorm.io/datatypes"
)

/*
  @Author : lanyulei
*/

// 字段数据
type Data struct {
	Id     int            `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`                 // 字段分组ID
	InfoId int            `gorm:"column:info_id; type:int(11);" json:"info_id" binding:"required"` // 对应的模型ID
	Status int            `gorm:"column:status; type:int(11); default:1" json:"status"`            // 0 没有状态，1 空闲，2 故障，3 待回收，4 正在使用
	Data   datatypes.JSON `gorm:"column:data; type:json" json:"data" binding:"required"`           // 数据
	models.BaseModel
}

func (Data) TableName() string {
	return "cmdb_resource_data"
}
