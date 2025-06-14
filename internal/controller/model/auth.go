package model

import "time"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Username            string    `json:"username"`
	AccessToken         string    `json:"access_token"`
	RefreshToken        string    `json:"refresh_token"`
	ExpiredToken        time.Time `json:"expired_token"`
	ExpiredRefreshToken time.Time `json:"expired_refresh_token"`
}
