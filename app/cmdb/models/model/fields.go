package model

import (
	"fiy/common/models"

	"gorm.io/datatypes"
)

/*
  @Author : lanyulei
*/

// 自定义字段
type Fields struct {
	Id              int            `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Identifies      string         `gorm:"column:identifies; type:varchar(128); uniqueIndex:idx_identifies_info_id;" json:"identifies" binding:"required"`                    // 字段唯一标识，英文的
	Name            string         `gorm:"column:name; type:varchar(128); uniqueIndex:idx_name_info_id;" json:"name" binding:"required"`                                      // 字段名称
	Type            string         `gorm:"column:type; type:varchar(45)" json:"type" binding:"required"`                                                                      // 字段类型(英文)。字符，数字...
	IsEdit          bool           `gorm:"column:is_edit; type:tinyint(1)" json:"is_edit"`                                                                                    // 是否可编辑
	IsUnique        bool           `gorm:"column:is_unique; type:tinyint(1)" json:"is_unique"`                                                                                // 是否唯一
	Required        bool           `gorm:"column:required; type:tinyint(1)" json:"required"`                                                                                  // 是否必填
	IsInternal      bool           `gorm:"column:is_internal; type:tinyint(1)" json:"is_internal"`                                                                            // 是否内置
	Prompt          string         `gorm:"column:prompt; type:varchar(1024)" json:"prompt"`                                                                                   // 用户提示
	Configuration   datatypes.JSON `gorm:"column:configuration; type:json" json:"configuration"`                                                                              // 字段配置相关的，例如数字的最大，最小，步长，单位等数据
	FieldGroupId    int            `gorm:"column:field_group_id; type:int(11)" json:"field_group_id" binding:"required"`                                                      // 字段分组ID
	InfoId          int            `gorm:"column:info_id; type:int(11); uniqueIndex:idx_name_info_id; uniqueIndex:idx_identifies_info_id;" json:"info_id" binding:"required"` // 对应的模型ID
	IsListDisplay   bool           `gorm:"column:is_list_display; type:tinyint(1); default:0" json:"is_list_display"`                                                         // 是否列表展示
	ListDisplaySort int            `gorm:"column:list_display_sort; type:int(11); default:999" json:"list_display_sort"`                                                      // 列表展示顺序
	models.BaseModel
}

func (Fields) TableName() string {
	return "cmdb_model_fields"
}
