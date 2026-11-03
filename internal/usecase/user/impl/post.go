package impl

import (
	"context"
	"fmt"
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/repository"
	repositoryModel "go-template/internal/repository/model"
	"go-template/pkg/event"
	"go-template/pkg/hash"
	"go-template/pkg/response"
	"net/http"
	"time"
)

func (u *Usecase) CreateUser(ctx context.Context, req controllerModel.CreateUserRequest) (res response.Response, err error) {
	nilfunc := repository.Repository{}
	if u.Repository == nilfunc {
		err = fmt.Errorf("repo is nil in usecase")
		return response.ResponseGenerate(http.StatusInternalServerError, err, nil, nil, nil), err
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return response.ResponseGenerate(http.StatusInternalServerError, response.ErrorGenerate(tag, err), nil, nil, nil), err
	}

	newUser := repositoryModel.User{
		Name:     req.Name,
		Username: req.Username,
		Password: hashedPassword,
	}

	created, err := u.Repository.DBRepository.UserRepository.CreateUserRepository(ctx, newUser)
	if err != nil {
		wrappedErr := fmt.Errorf("[CreateUser]-%s-%w", tag, err)
		return response.ResponseGenerate(http.StatusInternalServerError, wrappedErr, nil, nil, nil), wrappedErr
	}

	// Enqueue user.created event to Transactional Outbox
	eventID := fmt.Sprintf("user-created-%s-%d", created.Username, time.Now().UnixNano())
	_, cloudEventBytes, evtErr := event.NewCloudEvent(eventID, "user.created", "go-template.user", map[string]any{
		"id":         created.ID,
		"name":       created.Name,
		"username":   created.Username,
		"created_at": time.Now().UTC(),
	})
	if evtErr == nil && u.Repository.DBRepository.OutboxRepository != nil {
		_ = u.Repository.DBRepository.OutboxRepository.CreateOutboxEvent(ctx, repositoryModel.OutboxEvent{
			AggregateType: "user",
			AggregateID:   created.Username,
			Topic:         "user.events",
			Payload:       string(cloudEventBytes),
			Status:        repositoryModel.OutboxStatusPending,
		})
	}

	respData := controllerModel.User{
		ID:       created.ID,
		Name:     created.Name,
		Username: created.Username,
	}

	return response.ResponseGenerate(http.StatusCreated, nil, nil, respData, nil), nil
}

