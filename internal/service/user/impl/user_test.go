package impl

import (
	"context"
	"errors"
	controllerModel "go-template/internal/controller/model"
	"go-template/internal/outbound"
	outboundCache "go-template/internal/outbound/cache"
	"go-template/internal/repository"
	repositoryModel "go-template/internal/repository/model"
	"go-template/pkg/hash"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Definitions using testify/mock ---

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserRepository(ctx context.Context) ([]repositoryModel.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositoryModel.User), args.Error(1)
}

func (m *MockUserRepo) CreateUserRepository(ctx context.Context, user repositoryModel.User) (repositoryModel.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(repositoryModel.User), args.Error(1)
}

type MockEventQueueRepo struct {
	mock.Mock
}

func (m *MockEventQueueRepo) CreateEventQueue(ctx context.Context, event repositoryModel.EventQueue) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventQueueRepo) FetchPendingEvents(ctx context.Context, limit int) ([]repositoryModel.EventQueue, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]repositoryModel.EventQueue), args.Error(1)
}

func (m *MockEventQueueRepo) MarkProcessed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventQueueRepo) MarkFailed(ctx context.Context, id int64, retryCount, maxRetries int, lastErr string) error {
	args := m.Called(ctx, id, retryCount, maxRetries, lastErr)
	return args.Error(0)
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	// Directly execute the callback function simulating transaction commit
	return fn(ctx)
}

type MockUserCacheOutbound struct {
	mock.Mock
}

func (m *MockUserCacheOutbound) GetUsersCache(ctx context.Context) ([]controllerModel.User, bool, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).([]controllerModel.User), args.Bool(1), args.Error(2)
}

func (m *MockUserCacheOutbound) SetUsersCache(ctx context.Context, users []controllerModel.User) error {
	args := m.Called(ctx, users)
	return args.Error(0)
}

func (m *MockUserCacheOutbound) InvalidateUsersCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// --- Unit Test Scenarios: GetUsers ---

// Skenario 1: Cache Hit (Data ada di Redis)
// Ekspektasi: Langsung mengembalikan data cache, TIDAK query ke Database PostgreSQL.
func TestGetUsers_CacheHit(t *testing.T) {
	mockCache := new(MockUserCacheOutbound)
	mockRepo := new(MockUserRepo)

	cachedUsers := []controllerModel.User{
		{ID: 1, Name: "Budi Santoso", Username: "budi123"},
	}

	// Setup expectation: Redis hit = true
	mockCache.On("GetUsersCache", mock.Anything).Return(cachedUsers, true, nil).Once()

	svc := New(repository.Repository{
		UserRepository: mockRepo,
	}, outbound.Outbound{
		Cache: outboundCache.Outbound{User: mockCache},
	})

	res, err := svc.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.ResponseCode)
	assert.Equal(t, cachedUsers, res.Data)

	// Pastikan Database TIDAK pernah dipanggil saat cache hit
	mockRepo.AssertNotCalled(t, "GetUserRepository", mock.Anything)
	mockCache.AssertExpectations(t)
}

// Skenario 2: Cache Miss -> Query Database Sukses -> Simpan ke Redis
// Ekspektasi: Redis kosong -> query DB -> mapping data -> simpan ke Redis -> return 200 OK.
func TestGetUsers_CacheMiss_Success(t *testing.T) {
	mockCache := new(MockUserCacheOutbound)
	mockRepo := new(MockUserRepo)

	dbUsers := []repositoryModel.User{
		{ID: 1, Name: "Budi Santoso", Username: "budi123"},
		{ID: 2, Name: "Siti Rahma", Username: "siti456"},
	}

	// 1. Cache miss
	mockCache.On("GetUsersCache", mock.Anything).Return(nil, false, nil).Once()
	// 2. Query DB
	mockRepo.On("GetUserRepository", mock.Anything).Return(dbUsers, nil).Once()
	// 3. Simpan ke Cache
	mockCache.On("SetUsersCache", mock.Anything, mock.Anything).Return(nil).Once()

	svc := New(repository.Repository{
		UserRepository: mockRepo,
	}, outbound.Outbound{
		Cache: outboundCache.Outbound{User: mockCache},
	})

	res, err := svc.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.ResponseCode)

	users, ok := res.Data.([]controllerModel.User)
	assert.True(t, ok)
	assert.Len(t, users, 2)
	assert.Equal(t, "budi123", users[0].Username)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// Skenario 3: Cache Miss -> Query Database Gagal / Error
// Ekspektasi: Mengembalikan status 500 Internal Server Error beserta error message.
func TestGetUsers_DBError(t *testing.T) {
	mockCache := new(MockUserCacheOutbound)
	mockRepo := new(MockUserRepo)

	mockCache.On("GetUsersCache", mock.Anything).Return(nil, false, nil).Once()
	mockRepo.On("GetUserRepository", mock.Anything).Return(nil, errors.New("connection refused")).Once()

	svc := New(repository.Repository{
		UserRepository: mockRepo,
	}, outbound.Outbound{
		Cache: outboundCache.Outbound{User: mockCache},
	})

	res, err := svc.GetUsers(context.Background())

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.ResponseCode)
	assert.Nil(t, res.Data)

	mockRepo.AssertExpectations(t)
}

// --- Unit Test Scenarios: CreateUser ---

// Skenario 4: Create User Sukses
// Ekspektasi:
// 1. Password di-hash dengan bcrypt.
// 2. User baru disimpan ke database dalam transaksi.
// 3. Event "user.created" dicatat ke tabel EVENT_QUEUE.
// 4. Cache Redis di-invalidasi (dihapus agar sinkron).
// 5. Mengembalikan 201 Created.
func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockEventQueue := new(MockEventQueueRepo)
	mockTx := new(MockTxManager)
	mockCache := new(MockUserCacheOutbound)

	req := controllerModel.CreateUserRequest{
		Name:     "Andi Wijaya",
		Username: "andi123",
		Password: "SecretPassword123!",
	}

	createdModel := repositoryModel.User{
		ID:       10,
		Name:     req.Name,
		Username: req.Username,
	}

	// Expect UserRepository.CreateUserRepository dipanggil dengan password yang sudah di-hash
	mockRepo.On("CreateUserRepository", mock.Anything, mock.MatchedBy(func(u repositoryModel.User) bool {
		return u.Username == req.Username && hash.CheckPasswordHash(req.Password, u.Password)
	})).Return(createdModel, nil).Once()

	// Expect EventQueueRepository.CreateEventQueue dipanggil dengan topic "user.events"
	mockEventQueue.On("CreateEventQueue", mock.Anything, mock.MatchedBy(func(eq repositoryModel.EventQueue) bool {
		return eq.Topic == "user.events" && eq.AggregateID == req.Username
	})).Return(nil).Once()

	// Expect Invalidate cache
	mockCache.On("InvalidateUsersCache", mock.Anything).Return(nil).Once()

	svc := New(repository.Repository{
		TransactionManager:   mockTx,
		UserRepository:       mockRepo,
		EventQueueRepository: mockEventQueue,
	}, outbound.Outbound{
		Cache: outboundCache.Outbound{User: mockCache},
	})

	res, err := svc.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.ResponseCode)

	respUser, ok := res.Data.(controllerModel.User)
	assert.True(t, ok)
	assert.Equal(t, 10, respUser.ID)
	assert.Equal(t, req.Username, respUser.Username)

	mockRepo.AssertExpectations(t)
	mockEventQueue.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// Skenario 5: Create User Gagal karena Duplikat / Database Error
// Ekspektasi: Jika insert DB gagal, transaksi berhenti dan mengembalikan 500 error.
func TestCreateUser_DBError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockEventQueue := new(MockEventQueueRepo)
	mockTx := new(MockTxManager)

	req := controllerModel.CreateUserRequest{
		Name:     "Andi Wijaya",
		Username: "andi123",
		Password: "SecretPassword123!",
	}

	mockRepo.On("CreateUserRepository", mock.Anything, mock.Anything).
		Return(repositoryModel.User{}, errors.New("duplicate key violates unique constraint")).Once()

	svc := New(repository.Repository{
		TransactionManager:   mockTx,
		UserRepository:       mockRepo,
		EventQueueRepository: mockEventQueue,
	}, outbound.Outbound{})

	res, err := svc.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.ResponseCode)

	// EventQueue TIDAK boleh dipanggil jika insert user gagal
	mockEventQueue.AssertNotCalled(t, "CreateEventQueue", mock.Anything, mock.Anything)
}
