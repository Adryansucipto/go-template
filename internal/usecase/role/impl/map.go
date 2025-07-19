package impl

import (
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/repository/model"
	"time"
)

func MappingRoleToRequest(username string, request controllerModel.CreateRoleRequest) model.Role {
	return model.Role{
		Name:        request.Name,
		CreatedDate: time.Now(),
		CreatedBy:   username,
	}
}
