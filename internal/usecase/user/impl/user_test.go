package impl

import (
	"context"
	"errors"
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/repository"
	dbRepo "go-template/internal/repository/db"
	outboxRepo "go-template/internal/repository/db/outbox"
	userRepo "go-template/internal/repository/db/user"
	repositoryModel "go-template/internal/repository/model"
	"go-template/pkg/hash"
	"net/http"
	"testing"
)

// MockUserRepository implements userRepo.Repository for unit testing
type MockUserRepository struct {
	Users     []repositoryModel.User
	ShouldErr bool
}

func (m *MockUserRepository) GetUserRepository(ctx context.Context) ([]repositoryModel.User, error) {
	if m.ShouldErr {
		return nil, errors.New("database error")
	}
	return m.Users, nil
}

func (m *MockUserRepository) CreateUserRepository(ctx context.Context, u repositoryModel.User) (repositoryModel.User, error) {
	if m.ShouldErr {
		return repositoryModel.User{}, errors.New("failed to insert user")
	}
	u.ID = len(m.Users) + 1
	m.Users = append(m.Users, u)
	return u, nil
}

// MockOutboxRepository implements outboxRepo.Repository for unit testing
type MockOutboxRepository struct {
	Events []repositoryModel.OutboxEvent
}

func (m *MockOutboxRepository) CreateOutboxEvent(ctx context.Context, evt repositoryModel.OutboxEvent) error {
	m.Events = append(m.Events, evt)
	return nil
}

func (m *MockOutboxRepository) FetchPendingEvents(ctx context.Context, limit int) ([]repositoryModel.OutboxEvent, error) {
	return m.Events, nil
}

func (m *MockOutboxRepository) MarkProcessed(ctx context.Context, id int64) error {
	return nil
}

func (m *MockOutboxRepository) MarkFailed(ctx context.Context, id int64, retryCount int, maxRetries int, lastErr string) error {
	return nil
}

var (
	_ userRepo.Repository   = (*MockUserRepository)(nil)
	_ outboxRepo.Repository = (*MockOutboxRepository)(nil)
)

func TestGetUsers_Success(t *testing.T) {
	mockUser := &MockUserRepository{
		Users: []repositoryModel.User{
			{ID: 1, Name: "Budi", Username: "budi"},
			{ID: 2, Name: "Siti", Username: "siti"},
		},
	}

	uc := New(repository.Repository{
		DBRepository: dbRepo.Repository{
			UserRepository: mockUser,
		},
	})

	res, err := uc.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ResponseCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.ResponseCode)
	}

	users, ok := res.Data.([]controllerModel.User)
	if !ok || len(users) != 2 {
		t.Fatalf("expected 2 users, got %+v", res.Data)
	}
}

func TestCreateUser_SuccessAndEmitsOutboxEvent(t *testing.T) {
	mockUser := &MockUserRepository{}
	mockOutbox := &MockOutboxRepository{}

	uc := New(repository.Repository{
		DBRepository: dbRepo.Repository{
			UserRepository:   mockUser,
			OutboxRepository: mockOutbox,
		},
	})

	req := controllerModel.CreateUserRequest{
		Name:     "Andi Wijaya",
		Username: "andi123",
		Password: "PlainPassword123!",
	}

	res, err := uc.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ResponseCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.ResponseCode)
	}

	if len(mockUser.Users) != 1 {
		t.Fatalf("expected 1 user in repository, got %d", len(mockUser.Users))
	}
	if !hash.CheckPasswordHash(req.Password, mockUser.Users[0].Password) {
		t.Fatalf("expected stored password to be bcrypt hashed")
	}

	if len(mockOutbox.Events) != 1 {
		t.Fatalf("expected 1 outbox event emitted, got %d", len(mockOutbox.Events))
	}
	if mockOutbox.Events[0].Topic != "user.events" {
		t.Fatalf("expected outbox topic 'user.events', got %s", mockOutbox.Events[0].Topic)
	}
}

