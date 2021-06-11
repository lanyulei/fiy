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
	Id       int            `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`                          // 字段分组ID
	Uuid     string         `gorm:"column:uuid; type:varchar(45); unique;" json:"uuid" binding:"required"`    // 设备唯一ID
	InfoId   int            `gorm:"column:info_id; type:int(11); index;" json:"info_id" binding:"required"`   // 对应的模型ID
	InfoName string         `gorm:"column:info_name; type:varchar(128);" json:"info_name" binding:"required"` // 对应的模型名称
	Status   int            `gorm:"column:status; type:int(11); default:1" json:"status"`                     // 0 没有状态，1 空闲，2 故障，3 待回收，4 正在使用
	Data     datatypes.JSON `gorm:"column:data; type:json" json:"data" binding:"required"`                    // 数据
	models.BaseModel
}

func (Data) TableName() string {
	return "cmdb_resource_data"
}
