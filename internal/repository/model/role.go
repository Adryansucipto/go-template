package model

import "time"

type Role struct {
	ID           int        `gorm:"column:id"`
	Name         string     `gorm:"column:name"`
	CreatedDate  time.Time  `gorm:"column:created_date"`
	CreatedBy    string     `gorm:"column:created_by"`
	ModifiedDate *time.Time `gorm:"column:modified_date"`
	ModifiedBy   *string    `gorm:"column:modified_by"`
}

func (t *Role) TableName() string {
	return `"public"."MST_ROLE"`
}
