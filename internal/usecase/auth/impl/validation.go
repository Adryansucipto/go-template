package impl

import (
	"context"
	"errors"
	"fmt"
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/repository/model"
	"go-template/util"
	CONST "go-template/util/consts"
	"reflect"
	"strings"
)

func (u *Usecase) CheckActiveSession(ctx context.Context, username string) (int, error) {
	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSession(ctx, username)
	if err != nil {
		return 500, util.ErrorGenerate(
			tag,
			err,
		)
	}

	if !reflect.DeepEqual(record, model.Session{}) {
		return 400, util.ErrorGenerate(
			tag,
			errors.New("session still active, token not expired"),
		)
	}

	return 200, nil
}

func (u *Usecase) ValidateToken(ctx context.Context, token string) (string, util.Response, error) {
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return "", util.ResponseGenerate(
			400,
			CONST.ErrorMessages[400],
			"Error Token Not Found",
			nil,
			nil,
		), util.ErrorGenerate(tag, errors.New("error token not found"))
	}

	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByToken(ctx, strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			nil,
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}
	//check token habis?
	if reflect.DeepEqual(record, model.Session{}) {
		return "", util.ResponseGenerate(
			400,
			CONST.ErrorMessages[400],
			"Error Token Not Found",
			nil,
			nil,
		), errors.New("error token not found")
	}

	// check apakah ini token beneran?
	err = u.Config.JWT.VerifyToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			nil,
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}

	return strings.TrimPrefix(token, "Bearer "), util.ResponseGenerate(200, CONST.ErrorMessages[200], nil, nil, nil), nil
}

func (u *Usecase) AuthorizeUser(ctx context.Context, token string) (controllerModel.User, util.Response, error) {
	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByToken(ctx, strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return controllerModel.User{}, util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			nil,
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}

	recordUser, err := u.Repository.DBRepository.UserRepository.GetUserByIDRepository(ctx, record.Username)
	if err != nil {
		return controllerModel.User{}, util.ResponseGenerate(
			500,
			CONST.ErrorMessages[500],
			nil,
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}

	res := u.MappingUserRepoToClient(recordUser)

	return res, util.ResponseGenerate(200, CONST.ErrorMessages[200], nil, res, nil), nil
}

func (u *Usecase) ValidateRefreshToken(ctx context.Context, token string) (string, util.Response, error) {
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return "",
			util.ResponseGenerate(
				400,
				CONST.ErrorMessages[400],
				"Create Token Error",
				nil,
				nil,
			), util.ErrorGenerate(tag, errors.New("error token not found"))
	}
	tokenNoBearer := strings.Split(token, "Bearer ")
	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByRefreshToken(ctx, tokenNoBearer[1])
	if err != nil {
		return "",
			util.ResponseGenerate(
				500,
				CONST.ErrorMessages[500],
				nil,
				nil,
				nil,
			), util.ErrorGenerate(tag, err)
	}
	//cehck token habis?
	if reflect.DeepEqual(record, model.Session{}) {
		fmt.Println("token expired")
		return "",
			util.ResponseGenerate(
				400,
				CONST.ErrorMessages[400],
				"Error Token Not Found",
				nil,
				nil,
			), util.ErrorGenerate(tag, errors.New("error token not found"))
	}

	// check apakah ini token beneran?
	err = u.Config.JWT.VerifyRefreshToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", util.ResponseGenerate(
			400,
			CONST.ErrorMessages[400],
			"Create Token Error",
			nil,
			nil,
		), util.ErrorGenerate(tag, err)
	}

	return record.Username, util.ResponseGenerate(200, CONST.ErrorMessages[200], nil, nil, nil), nil
}
