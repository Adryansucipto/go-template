package impl

import (
	"context"
	"fmt"
	"go-template/internal/controller/model"
	"go-template/internal/repository"
)

func (u *Usecase) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User

	nilfunc := repository.Repository{}
	if u.Repository == nilfunc {
		return nil, fmt.Errorf("repo is nil in usecase")
	}

	records, err := u.Repository.DBRepository.UserRepository.GetUserRepository(ctx)
	if err != nil {
		return []model.User{}, fmt.Errorf("[GetUsers]-%s-%+v", tag, err.Error())
	}

	//mapping

	for _, v := range records {
		var user model.User

		user.ID = v.ID
		user.Name = v.Name

		users = append(users, user)
	}

	return users, nil
}
