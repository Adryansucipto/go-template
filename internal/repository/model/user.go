package model

type User struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (t *User) TableName() string {
	return `"public"."MST_USER"`
}
