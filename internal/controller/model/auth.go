package model

import (
	"go-template/pkg/validator"
	"time"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r AuthRequest) Validate() error {
	return validator.RequireNonEmpty(map[string]string{
		"username": r.Username,
		"password": r.Password,
	})
}

type AuthResponse struct {
	Username            string    `json:"username"`
	AccessToken         string    `json:"access_token"`
	RefreshToken        string    `json:"refresh_token"`
	ExpiredToken        time.Time `json:"expired_token"`
	ExpiredRefreshToken time.Time `json:"expired_refresh_token"`
}
