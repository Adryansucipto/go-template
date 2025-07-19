package impl

import (
	"context"
	"errors"
	controllerModel "go-template/internal/controller/model"
	repositoryModel "go-template/internal/repository/model"
	"go-template/util"
	CONST "go-template/util/consts"
	"time"
)

func (u *Usecase) LoginHandler(ctx context.Context, request controllerModel.AuthRequest) (response util.Response, err error) {

	getRecord, err := u.Repository.DBRepository.AuthRepository.GetUsernamePassword(ctx, request.Username)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				nil,
				nil,
				nil,
			), util.ErrorGenerate(
				tag,
				err)
	}

	// validate
	if getRecord.Password != request.Password {
		return util.ResponseGenerate(
				400,
				CONST.ErrorMessages[400],
				errors.New("error not found"),
				nil,
				nil,
			), util.ErrorGenerate(
				tag,
				errors.New("error not found"))
	}

	//check session exists
	code, err := u.CheckActiveSession(ctx, request.Username)
	if err != nil {
		return util.ResponseGenerate(
				code,
				CONST.ErrorMessages[code],
				"Session Still Active",
				nil,
				nil,
			), util.ErrorGenerate(
				tag,
				err)
	}

	token, expiredToken, err := u.Config.JWT.CreateToken(request.Username)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				"Create Token Error",
				nil,
				nil,
			), util.ErrorGenerate(
				tag,
				err)
	}

	refreshToken, expiredRefreshToken, err := u.Config.JWT.CreateRefreshToken(request.Username)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				"Create Token Error",
				nil,
				nil), util.ErrorGenerate(
				tag,
				err)
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
	err = u.Repository.DBRepository.AuthRepository.DeleteSession(ctx, request.Username)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				"Create Token Error",
				nil,
				nil,
			), util.ErrorGenerate(
				tag,
				err)
	}

	err = u.Repository.DBRepository.AuthRepository.CreateSession(ctx, sessionRepo)
	if err != nil {
		return util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			util.ErrorGenerate(
				tag,
				err),
			nil,
			nil,
		), err
	}

	response = util.ResponseGenerate(200, CONST.ErrorMessages[200], nil, data, nil)

	return response, err
}

func (u *Usecase) RefreshNewToken(ctx context.Context, username string) (response util.Response, err error) {

	// Create new Access_Token
	token, expiredToken, err := u.Config.JWT.CreateToken(username)
	if err != nil {
		return util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			"Create Token Error",
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}

	refreshToken, expiredRefreshToken, err := u.Config.JWT.CreateRefreshToken(username)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				"Create Token Error",
				nil,
				nil),
			util.ErrorGenerate(tag, err)
	}

	var data controllerModel.AuthResponse
	data.Username = username
	data.AccessToken = token
	data.RefreshToken = refreshToken
	data.ExpiredToken = expiredToken
	data.ExpiredRefreshToken = expiredRefreshToken

	// Insert To Database
	var sessionRepo repositoryModel.Session
	sessionRepo.Username = username
	sessionRepo.AccessToken = token
	sessionRepo.RefreshToken = refreshToken
	sessionRepo.CreatedDate = time.Now()
	sessionRepo.ExpiredSession = expiredToken
	sessionRepo.ExpiredRefresh = expiredRefreshToken
	err = u.Repository.DBRepository.AuthRepository.DeleteSession(ctx, username)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				"Create Token Error",
				nil,
				nil,
			),
			util.ErrorGenerate(tag, err)
	}

	err = u.Repository.DBRepository.AuthRepository.CreateSession(ctx, sessionRepo)
	if err != nil {
		return util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				nil,
				nil,
				nil,
			),
			util.ErrorGenerate(tag, err)
	}

	response = util.ResponseGenerate(200, CONST.ErrorMessages[200], nil, data, nil)
	return response, nil
}
