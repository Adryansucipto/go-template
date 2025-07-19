package impl

import (
	"context"
	"go-template/internal/controller/model"
	"go-template/util"
	CONST "go-template/util/consts"
	"strings"
)

func (u Usecase) CreateRole(ctx context.Context, request model.CreateRoleRequest, username string) (util.Response, error) {

	// create roles
	// Mapping Roles to Repository

	requestRepository := MappingRoleToRequest(username, request)
	if err := u.Repository.DBRepository.RoleRepository.CreateRole(ctx, requestRepository); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return util.ResponseGenerate(
				400,
				CONST.ErrorMessages[400],
				"Role Name Already Exists",
				nil,
				nil,
			), util.ErrorGenerate(tag, err)
		}

		return util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			nil,
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}

	return util.ResponseGenerate(
		200,
		CONST.ErrorMessages[200],
		nil,
		nil,
		nil,
	), nil
}
