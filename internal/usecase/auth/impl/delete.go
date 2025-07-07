package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/util"
)

func (u *Usecase) DeleteSession(ctx context.Context, record repoModel.Session) util.Response {

	err := u.Repository.DBRepository.AuthRepository.DeleteSession(ctx, record.Username)
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
