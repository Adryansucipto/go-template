package impl

import (
	authRepository "go-template/internal/repository/auth"
	"go-template/pkg/resource/model"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

var (
	tag = "[AuthRepository]"
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
