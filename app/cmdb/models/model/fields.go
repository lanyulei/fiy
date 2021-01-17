package model

import (
	"encoding/json"
	"fiy/common/models"
)

/*
  @Author : lanyulei
*/

// 自定义字段
type Fields struct {
	Id            int             `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Identifies    string          `gorm:"column:identifies; type:varchar(128)" json:"identifies" binding:"required"` // 字段唯一标识，英文的
	Name          string          `gorm:"column:name; type:varchar(128)" json:"name" binding:"required"`             // 字段名称
	Type          string          `gorm:"column:type; type:varchar(45)" json:"type"`                                 // 字段类型。字符，数字...
	IsEdit        bool            `gorm:"column:is_edit; type:tinyint(1)" json:"is_edit"`                            // 是否可编辑
	IsUnique      bool            `gorm:"column:is_unique; type:tinyint(1)" json:"is_unique"`                        // 是否唯一
	IsInternal    bool            `gorm:"column:is_internal; type:tinyint(1)" json:"is_internal"`                    // 是否内置
	Regular       string          `gorm:"column:regular; type:varchar(1024)" json:"regular"`                         // 正则表达式
	Prompt        string          `gorm:"column:prompt; type:varchar(1024)" json:"prompt"`                           // 用户提示
	Configuration json.RawMessage `gorm:"column:configuration; type:json" json:"configuration"`                      // 字段配置相关的，例如数字的最大，最小，步长，单位等数据
	FieldGroupId  int             `gorm:"column:field_group_id; type:int(11)" json:"field_group_id"`                 // 字段分组ID
	InfoId        int             `gorm:"column:info_id; type:int(11)" json:"info_id"`                               // 对应的模型ID
	models.BaseModel
}

func (Fields) TableName() string {
	return "cmdb_model_fields"
}
