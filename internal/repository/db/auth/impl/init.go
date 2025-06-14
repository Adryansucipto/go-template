package impl

import (
	authRepository "go-template/internal/repository/db/auth"
	"go-template/pkg/resource/model"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Repository struct {
	dig.In

	DbGorm *gorm.DB
}

func New(db model.Database) authRepository.Repository {
	return &Repository{
		DbGorm: db.PostgreSql,
	}
}
