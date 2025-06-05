package impl

import (
	"context"
	"go-template/internal/controller/model"
)

func (u *Usecase) GetUsers(ctx context.Context) ([]model.User, error) {
	var users = []model.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	return users, nil
}
