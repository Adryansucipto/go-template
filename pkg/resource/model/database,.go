package model

import "gorm.io/gorm"

type Database struct {
	PostgreSql *gorm.DB
}
