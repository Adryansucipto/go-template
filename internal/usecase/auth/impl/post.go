package impl

import (
	"context"
	"errors"
	"go-template/internal/controller/model"
	"go-template/util"
)

func (u *Usecase) LoginHandler(ctx context.Context, request model.AuthRequest) (response util.Response, err error) {

	getRecord, err := u.Repository.DBRepository.AuthRepository.GetUsernamePassword(ctx, request.Username)
	if err != nil {
		return util.ResponseGenerate(500, errors.New("error not found"), nil, nil), err

	}

	// validate
	if getRecord.Password != request.Password {
		return util.ResponseGenerate(500, errors.New("error not found"), nil, nil), err
	}

	// check session based on outbound - karena login jadi check by DB
	// kalau action lain ada token baru check by func

	token, expiredToken, err := u.Config.JWT.CreateToken(request.Username)
	if err != nil {
		return util.ResponseGenerate(500, errors.New("error not found"), nil, nil), err
	}

	refreshToken, expiredRefreshToken, err := u.Config.JWT.CreateRefreshToken(request.Username)
	if err != nil {
		return util.ResponseGenerate(500, errors.New("error not found"), nil, nil), err
	}

	var data model.AuthResponse
	data.Username = request.Username
	data.AccessToken = token
	data.RefreshToken = refreshToken
	data.ExpiredToken = expiredToken
	data.ExpiredRefreshToken = expiredRefreshToken

	response = util.ResponseGenerate(200, nil, data, nil)

	return response, err
}
