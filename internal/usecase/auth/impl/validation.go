package impl

import (
	"context"
	"errors"
	"fmt"
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/repository/model"
	"go-template/util"
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

func (u *Usecase) ValidateToken(ctx context.Context, token string) (string, util.Response) {
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return "", util.ResponseGenerate(
			400,
			errors.New("error token not found"),
			nil,
			nil,
			nil,
		)
	}

	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByToken(ctx, strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", util.ResponseGenerate(
			500,
			err,
			nil,
			nil,
			nil,
		)
	}
	//check token habis?
	if reflect.DeepEqual(record, model.Session{}) {
		return "", util.ResponseGenerate(
			400,
			errors.New("error token not found"),
			nil,
			nil,
			nil,
		)
	}

	// check apakah ini token beneran?
	err = u.Config.JWT.VerifyToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", util.ResponseGenerate(
			400,
			err,
			nil,
			nil,
			nil,
		)
	}

	return strings.TrimPrefix(token, "Bearer "), util.ResponseGenerate(200, nil, nil, nil, nil)
}

func (u *Usecase) AuthorizeUser(ctx context.Context, token string) (controllerModel.User, util.Response) {
	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByToken(ctx, strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return controllerModel.User{}, util.ResponseGenerate(
			500,
			err,
			nil,
			nil,
			nil,
		)
	}

	recordUser, err := u.Repository.DBRepository.UserRepository.GetUserByIDRepository(ctx, record.Username)
	if err != nil {
		return controllerModel.User{}, util.ResponseGenerate(
			500,
			err,
			nil,
			nil,
			nil,
		)
	}

	res := u.MappingUserRepoToClient(recordUser)

	return res, util.ResponseGenerate(200, nil, nil, nil, nil)
}

func (u *Usecase) ValidateRefreshToken(ctx context.Context, token string) (string, util.Response) {
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return "",
			util.ResponseGenerate(
				400,
				errors.New("error token not found"),
				nil,
				nil,
				nil,
			)
	}
	tokenNoBearer := strings.Split(token, "Bearer ")
	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByRefreshToken(ctx, tokenNoBearer[1])
	if err != nil {
		return "",
			util.ResponseGenerate(
				500,
				err,
				nil,
				nil,
				nil,
			)
	}
	//cehck token habis?
	if reflect.DeepEqual(record, model.Session{}) {
		fmt.Println("token expired")
		return "",
			util.ResponseGenerate(
				400,
				errors.New("error token not found"),
				nil,
				nil,
				nil,
			)
	}

	// check apakah ini token beneran?
	err = u.Config.JWT.VerifyRefreshToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", util.ResponseGenerate(
			400,
			err,
			nil,
			nil,
			nil,
		)
	}

	return record.Username, util.ResponseGenerate(200, nil, nil, nil, nil)
}
