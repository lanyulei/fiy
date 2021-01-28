package models

import (
	"time"
)

type BaseModel struct {
	CreatedAt time.Time  `gorm:"type: datetime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"type: datetime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index; type: datetime" json:"-"`
}
