package model

import "go-template/pkg/validator"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username,omitempty"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Validate() error {
	return validator.RequireNonEmpty(map[string]string{
		"name":     r.Name,
		"username": r.Username,
		"password": r.Password,
	})
}
