package impl

import (
	"context"
	"errors"
	"fmt"
	controllerModel "go-template/internal/controller/model"
	repositoryModel "go-template/internal/repository/model"
	"go-template/pkg/event"
	"go-template/pkg/hash"
	"go-template/pkg/response"
	"strings"
	"time"
)

func (u *Usecase) LoginHandler(ctx context.Context, request controllerModel.AuthRequest) (res response.Response) {
	getRecord, err := u.Repository.DBRepository.AuthRepository.GetUsernamePassword(ctx, request.Username)
	if err != nil {
		return response.ResponseGenerate(
			401,
			response.ErrorGenerate(tag, errors.New("invalid username or password")),
			nil,
			nil,
			nil,
		)
	}

	// Support both bcrypt hashed passwords and legacy plaintext during migration
	passwordValid := false
	if strings.HasPrefix(getRecord.Password, "$2a$") || strings.HasPrefix(getRecord.Password, "$2b$") {
		passwordValid = hash.CheckPasswordHash(request.Password, getRecord.Password)
	} else {
		passwordValid = getRecord.Password == request.Password
	}

	if !passwordValid {
		return response.ResponseGenerate(
			401,
			response.ErrorGenerate(tag, errors.New("invalid username or password")),
			nil,
			nil,
			nil,
		)
	}

	code, err := u.CheckTokenExpired(ctx, "", request.Username)
	if err != nil {
		return response.ResponseGenerate(
			code,
			response.ErrorGenerate(tag, err),
			nil,
			nil,
			nil,
		)
	}

	token, expiredToken, err := u.Config.JWT.CreateToken(request.Username)
	if err != nil {
		return response.ResponseGenerate(
			500,
			response.ErrorGenerate(tag, err),
			nil,
			nil,
			nil,
		)
	}

	refreshToken, expiredRefreshToken, err := u.Config.JWT.CreateRefreshToken(request.Username)
	if err != nil {
		return response.ResponseGenerate(
			500,
			response.ErrorGenerate(tag, err),
			nil,
			nil,
			nil,
		)
	}

	var data controllerModel.AuthResponse
	data.Username = request.Username
	data.AccessToken = token
	data.RefreshToken = refreshToken
	data.ExpiredToken = expiredToken
	data.ExpiredRefreshToken = expiredRefreshToken

	var sessionRepo repositoryModel.Session
	sessionRepo.Username = request.Username
	sessionRepo.AccessToken = token
	sessionRepo.RefreshToken = refreshToken
	sessionRepo.CreatedDate = time.Now()
	sessionRepo.ExpiredSession = expiredToken
	sessionRepo.ExpiredRefresh = expiredRefreshToken

	if err = u.Repository.DBRepository.AuthRepository.DeleteSession(ctx, request.Username); err != nil {
		return response.ResponseGenerate(500, err, nil, nil, nil)
	}

	if err = u.Repository.DBRepository.AuthRepository.CreateSession(ctx, sessionRepo); err != nil {
		return response.ResponseGenerate(500, err, nil, nil, nil)
	}

	// Enqueue UserLoginEvent to Transactional Outbox table
	eventID := fmt.Sprintf("user-login-%s-%d", request.Username, time.Now().UnixNano())
	_, cloudEventBytes, err := event.NewCloudEvent(eventID, "user.logged_in", "go-template.auth", map[string]any{
		"username":  request.Username,
		"logged_at": time.Now().UTC(),
	})
	if err == nil {
		_ = u.Repository.DBRepository.OutboxRepository.CreateOutboxEvent(ctx, repositoryModel.OutboxEvent{
			AggregateType: "user",
			AggregateID:   request.Username,
			Topic:         "user.events",
			Payload:       string(cloudEventBytes),
			Status:        repositoryModel.OutboxStatusPending,
		})
	}

	return response.ResponseGenerate(200, nil, nil, data, nil)
}
