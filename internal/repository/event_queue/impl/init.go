package impl

import (
	eventQueueRepo "go-template/internal/repository/event_queue"
	"go-template/pkg/resource/model"

	"gorm.io/gorm"
)

var (
	tag = "[EventQueueRepository]"
)

type Repository struct {
	DbGorm *gorm.DB
}

func New(db model.Database) eventQueueRepo.Repository {
	return &Repository{
		DbGorm: db.PostgreSql,
	}
}

