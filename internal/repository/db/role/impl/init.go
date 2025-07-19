package impl

import (
	roleRepository "go-template/internal/repository/db/role"
	"go-template/pkg/resource/model"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

var (
	tag = "[RoleRepository]"
)

type Repository struct {
	dig.In

	DbGorm *gorm.DB
}

func New(db model.Database) roleRepository.Repository {
	return &Repository{
		DbGorm: db.PostgreSql,
	}
}
