package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCnfg struct {
	Secret string `yaml:"secret"`
}

func (jwts *JwtCnfg) CreateToken(username string) (string, time.Time, error) {
	expired := time.Now().Add(time.Minute * 20)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"type":     "access",
			"exp":      expired.Unix(),
		})

	tokenString, err := token.SignedString([]byte(jwts.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expired, nil
}

func (jwts *JwtCnfg) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwts.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "access" {
		return nil, fmt.Errorf("not an access token")
	}

	return claims, nil
}

func (jwts *JwtCnfg) CreateRefreshToken(username string) (string, time.Time, error) {
	expired := time.Now().Add(time.Hour * 24 * 7)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"type":     "refresh",
		"exp":      expired.Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwts.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expired, nil
}

func (jwts *JwtCnfg) VerifyRefreshToken(refreshToken string) (string, time.Time, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwts.Secret), nil
	})

	if err != nil {
		return "", time.Time{}, err
	}
	if !token.Valid {
		return "", time.Time{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", time.Time{}, errors.New("invalid refresh token type")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", time.Time{}, errors.New("invalid username claim")
	}

	newAccessToken, expiredToken, err := jwts.CreateToken(username)
	if err != nil {
		return "", time.Time{}, err
	}

	return newAccessToken, expiredToken, nil
}
