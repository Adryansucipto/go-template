package model

type User struct {
	ID       int    `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (t *User) TableName() string {
	return `"public"."MST_USER"`
}
