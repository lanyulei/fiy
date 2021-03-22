package resource

import (
	"database/sql"
	"fiy/common/models"

	"gorm.io/datatypes"
)

/*
  @Author : lanyulei
*/

// 云资源同步管理
type CloudDiscovery struct {
	Id             int            `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name           string         `gorm:"column:name; type:varchar(128); unique; uniqueIndex:idx_unique;" json:"name" binding:"required"`       // 任务名称
	ResourceModel  int            `gorm:"column:resource_model; type:int(11);uniqueIndex:idx_unique;" json:"resource_model" binding:"required"` // 资源模型
	ResourceType   int            `gorm:"column:resource_type; type:int(11);" json:"resource_type" binding:"required"`                          // 资源类型
	CloudAccount   int            `gorm:"column:cloud_account; type:int(11);uniqueIndex:idx_unique;" json:"cloud_account" binding:"required"`   // 云账户
	Region         datatypes.JSON `gorm:"column:region; type:json;" json:"region" binding:"required"`                                           // 区域
	Status         bool           `gorm:"column:status; type:tinyint(1); default:1" json:"status"`                                              // 任务状态
	FieldMap       datatypes.JSON `gorm:"column:field_map; type:json;" json:"field_map"`                                                        // 字段映射 [{"source": "hostname", "target": "hostname"}]
	LastSyncStatus string         `gorm:"column:last_sync_status; type:varchar(45);" json:"last_sync_status"`                                   // 最近同步状态
	LastSyncTime   sql.NullTime   `gorm:"column:last_sync_time; type:datetime;" json:"last_sync_time"`                                          // 最近同步状态
	Creator        int            `gorm:"column:creator; type:int(11);" json:"creator"`                                                         // 创建者
	Modifier       int            `gorm:"column:modifier; type:int(11);" json:"modifier"`                                                       // 修改者
	Remarks        string         `gorm:"column:remarks; type:varchar(1024);" json:"remarks"`                                                   // 备注
	models.BaseModel
}

func (CloudDiscovery) TableName() string {
	return "cmdb_resource_cloud_discovery"
}
