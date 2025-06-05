package impl

import (
	userRepository "go-template/internal/repository/db/user"

	"gorm.io/gorm"
)

type Repository struct {
	DbGorm *gorm.DB
}

func New() userRepository.Repository {
	return &Repository{}
}
