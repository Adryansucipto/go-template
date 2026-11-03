package impl

import (
	outboxRepo "go-template/internal/repository/db/outbox"
	"go-template/pkg/resource/model"

	"gorm.io/gorm"
)

var (
	tag = "[OutboxDBRepository]"
)

type Repository struct {
	DbGorm *gorm.DB
}

func New(db model.Database) outboxRepo.Repository {
	return &Repository{
		DbGorm: db.PostgreSql,
	}
}

