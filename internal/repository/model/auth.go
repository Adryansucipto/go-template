package model

import "time"

type Session struct {
	Username       string    `gorm:"column:username"`
	AccessToken    string    `gorm:"column:access_token"`
	RefreshToken   string    `gorm:"column:refresh_token"`
	CreatedDate    time.Time `gorm:"column:created_date"`
	ExpiredSession time.Time `gorm:"column:expired_session"`
	ExpiredRefresh time.Time `gorm:"column:expired_refresh"`
}

func (t *Session) TableName() string {
	return `"public"."SESSION"`
}
