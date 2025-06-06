package impl

import (
	userRepository "go-template/internal/repository/db/user"
	"go-template/pkg/resource/model"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Repository struct {
	dig.In

	DbGorm *gorm.DB
}

func New(db model.Database) userRepository.Repository {
	return &Repository{
		DbGorm: db.PostgreSql,
	}
}
