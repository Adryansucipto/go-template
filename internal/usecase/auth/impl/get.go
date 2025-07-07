package impl

import (
	"context"
	"errors"
	"fmt"
	"go-template/internal/repository/model"
	"go-template/util"
	"reflect"
	"strings"
)

func (u *Usecase) GetRecordByToken(ctx context.Context, token string) (model.Session, util.Response) {
	if token == "" {
		return model.Session{},
			util.ResponseGenerate(
				400,
				errors.New("error token not found"),
				nil,
				nil,
				nil,
			)
	}
	tokenNoBearer := strings.Split(token, "Bearer ")
	record, err := u.Repository.DBRepository.AuthRepository.GetActiveSessionByToken(ctx, tokenNoBearer[1])
	if err != nil {
		return model.Session{},
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
		return model.Session{},
			util.ResponseGenerate(
				400,
				errors.New("error token not found"),
				nil,
				nil,
				nil,
			)
	}

	return record, util.ResponseGenerate(200, nil, nil, nil, nil)
}
