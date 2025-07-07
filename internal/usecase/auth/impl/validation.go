package impl

import (
	"context"
	"errors"
	"go-template/internal/repository/model"
	"go-template/util"
	"reflect"
)

func (u *Usecase) CheckTokenExpired(ctx context.Context, token, username string) (int, error) {
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
