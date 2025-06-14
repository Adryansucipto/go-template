package model

type Session struct {
	Username     string `gorm:"column:username"`
	AccessToken  string `gorm:"access_token"`
	RefreshToken string `gorm:"refresh_token"`
	CreatedDate  string `gorm:"created_date"`
}

func (t *Session) TableName() string {
	return `"public"."SESSION"`
}
