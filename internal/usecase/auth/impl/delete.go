package impl

import (
	"context"
	"go-template/util"
)

func (u *Usecase) DeleteSession(ctx context.Context, username string) util.Response {

	err := u.Repository.DBRepository.AuthRepository.DeleteSession(ctx, username)
	if err != nil {
		return util.ResponseGenerate(
			500,
			err,
			nil,
			nil,
			nil,
		)
	}
	return util.ResponseGenerate(
		200,
		nil,
		nil,
		nil,
		nil,
	)
}
