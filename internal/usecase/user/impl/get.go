package impl

import (
	"context"
	"fmt"
	"go-template/internal/controller/model"
	"go-template/internal/repository"
	"go-template/util"
)

func (u *Usecase) GetUsers(ctx context.Context) (res util.Response, err error) {
	var users []model.User

	nilfunc := repository.Repository{}
	if u.Repository == nilfunc {
		return util.ResponseGenerate(500, fmt.Errorf("repo is nil in usecase"), nil, nil), fmt.Errorf("repo is nil in usecase")
	}

	records, err := u.Repository.DBRepository.UserRepository.GetUserRepository(ctx)
	if err != nil {
		return util.ResponseGenerate(500, fmt.Errorf("[GetUsers]-%s-%+v", tag, err.Error()), nil, nil), fmt.Errorf("[GetUsers]-%s-%+v", tag, err.Error())
	}

	//mapping

	for _, v := range records {
		var user model.User

		user.ID = v.ID
		user.Name = v.Name

		users = append(users, user)
	}

	return util.ResponseGenerate(200, nil, users, nil), nil
}
