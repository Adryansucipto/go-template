package impl

import (
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/repository/model"
)

func (u *Usecase) MappingUserRepoToClient(record model.User) controllerModel.User {
	return controllerModel.User{
		ID:       record.ID,
		Name:     record.Name,
		Username: record.Username,
	}
}
