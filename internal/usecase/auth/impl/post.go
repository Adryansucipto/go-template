package impl

import (
	"context"
	"errors"
	"fmt"
	controllerModel "go-template/internal/controller/model"
	repositoryModel "go-template/internal/repository/model"
	"go-template/util"
	"time"
)

func (u *Usecase) LoginHandler(ctx context.Context, request controllerModel.AuthRequest) (response util.Response) {

	getRecord, err := u.Repository.DBRepository.AuthRepository.GetUsernamePassword(ctx, request.Username)
	if err != nil {
		return util.ResponseGenerate(
			500,
			util.ErrorGenerate(
				tag,
				err),
			nil,
			nil,
			nil,
		)
	}

	// validate
	if getRecord.Password != request.Password {
		return util.ResponseGenerate(
			500,
			util.ErrorGenerate(
				tag,
				errors.New("error not found")),
			nil,
			nil,
			nil,
		)
	}

	// check session based on outbound - karena login jadi check by DB
	// kalau action lain ada token baru check by func

	//check session exists
	active := u.Repository.DBRepository.AuthRepository.GetActiveSession(ctx, request.Username)
	fmt.Println("is active: ", active)
	if active {
		return util.ResponseGenerate(
			400,
			util.ErrorGenerate(
				tag,
				errors.New("error on Session"),
			),
			[]string{"Account is on session"},
			nil,
			nil,
		)
	} else {
		// delete session
	}

	token, expiredToken, err := u.Config.JWT.CreateToken(request.Username)
	if err != nil {
		return util.ResponseGenerate(
			500,
			util.ErrorGenerate(
				tag,
				err),
			nil,
			nil,
			nil,
		)
	}

	refreshToken, expiredRefreshToken, err := u.Config.JWT.CreateRefreshToken(request.Username)
	if err != nil {
		return util.ResponseGenerate(
			500,
			util.ErrorGenerate(
				tag,
				err),
			nil,
			nil,
			nil)
	}

	var data controllerModel.AuthResponse
	data.Username = request.Username
	data.AccessToken = token
	data.RefreshToken = refreshToken
	data.ExpiredToken = expiredToken
	data.ExpiredRefreshToken = expiredRefreshToken

	// Insert To Database
	var sessionRepo repositoryModel.Session
	sessionRepo.Username = request.Username
	sessionRepo.AccessToken = token
	sessionRepo.RefreshToken = refreshToken
	sessionRepo.CreatedDate = time.Now()
	sessionRepo.ExpiredSession = expiredToken
	sessionRepo.ExpiredRefresh = expiredRefreshToken
	err = u.Repository.DBRepository.AuthRepository.CreateSession(ctx, sessionRepo)
	if err != nil {
		return util.ResponseGenerate(
			500,
			err,
			nil,
			nil,
			nil,
		)
	}

	response = util.ResponseGenerate(200, nil, nil, data, nil)

	return response
}
