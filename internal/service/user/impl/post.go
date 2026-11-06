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

func (s *Service) CreateUser(ctx context.Context, req controllerModel.CreateUserRequest) (res response.Response, err error) {
	nilRepo := repository.Repository{}
	if s.Repository == nilRepo {
		err = fmt.Errorf("repo is nil in service")
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

	var created repositoryModel.User

	runInTx := func(txCtx context.Context) error {
		var txErr error
		created, txErr = s.Repository.UserRepository.CreateUserRepository(txCtx, newUser)
		if txErr != nil {
			return fmt.Errorf("[CreateUser]-%s-%w", tag, txErr)
		}

		if s.Repository.EventQueueRepository != nil {
			eventID := fmt.Sprintf("user-created-%s-%d", created.Username, time.Now().UnixNano())
			_, cloudEventBytes, evtErr := event.NewCloudEvent(eventID, "user.created", "go-template.user", map[string]any{
				"id":         created.ID,
				"name":       created.Name,
				"username":   created.Username,
				"created_at": time.Now().UTC(),
			})
			if evtErr != nil {
				return evtErr
			}

			if txErr := s.Repository.EventQueueRepository.CreateEventQueue(txCtx, repositoryModel.EventQueue{
				AggregateType: "user",
				AggregateID:   created.Username,
				Topic:         "user.events",
				Payload:       string(cloudEventBytes),
				Status:        repositoryModel.EventQueueStatusPending,
			}); txErr != nil {
				return txErr
			}
		}

		return nil
	}

	// Execute User Creation + EventQueue Event insertion atomically in a single DB Transaction
	if s.Repository.TransactionManager != nil {
		err = s.Repository.TransactionManager.WithTransaction(ctx, runInTx)
	} else {
		err = runInTx(ctx)
	}

	if err != nil {
		return response.ResponseGenerate(http.StatusInternalServerError, err, nil, nil, nil), err
	}

	// Invalidate Redis cache via Outbound layer after successful transaction commit
	if s.Outbound.Cache.User != nil {
		_ = s.Outbound.Cache.User.InvalidateUsersCache(ctx)
	}

	respData := controllerModel.User{
		ID:       created.ID,
		Name:     created.Name,
		Username: created.Username,
	}

	return response.ResponseGenerate(http.StatusCreated, nil, nil, respData, nil), nil
}
