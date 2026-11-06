package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-template/internal/controller/model"
	"go-template/pkg/response"
	"go-template/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserService implements userService.Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUsers(ctx context.Context) (response.Response, error) {
	args := m.Called(ctx)
	return args.Get(0).(response.Response), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, req model.CreateUserRequest) (response.Response, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(response.Response), args.Error(1)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()
	return e
}

// Skenario 1: GET /v1/users Sukses
// Ekspektasi: Service mengembalikan data -> Handler mengembalikan HTTP Status 200 OK.
func TestGetUsers_Success(t *testing.T) {
	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	rec := httptest.NewRecorder()
	eCtx := e.NewContext(req, rec)

	mockSvc := new(MockUserService)
	expectedUsers := []model.User{
		{ID: 1, Name: "Budi Santoso", Username: "budi123"},
	}

	mockSvc.On("GetUsers", mock.Anything).Return(
		response.ResponseGenerate(http.StatusOK, nil, nil, expectedUsers, nil),
		nil,
	).Once()

	ctrl := Controller{
		User:   mockSvc,
		Logger: zap.NewNop(),
	}

	err := ctrl.GetUsers(eCtx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var respBody response.FinalResponse
	err = json.Unmarshal(rec.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, 200, respBody.ResponseCode)

	mockSvc.AssertExpectations(t)
}

// Skenario 2: GET /v1/users Service Error
// Ekspektasi: Service mengembalikan error -> Handler mengembalikan HTTP Status 500.
func TestGetUsers_ServiceError(t *testing.T) {
	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	rec := httptest.NewRecorder()
	eCtx := e.NewContext(req, rec)

	mockSvc := new(MockUserService)
	mockSvc.On("GetUsers", mock.Anything).Return(
		response.ResponseGenerate(http.StatusInternalServerError, errors.New("db down"), nil, nil, nil),
		errors.New("db down"),
	).Once()

	ctrl := Controller{
		User:   mockSvc,
		Logger: zap.NewNop(),
	}

	err := ctrl.GetUsers(eCtx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockSvc.AssertExpectations(t)
}

// Skenario 3: POST /v1/register Sukses
// Ekspektasi: Payload valid -> Service membuat user -> Handler mengembalikan HTTP Status 201 Created.
func TestCreateUser_Success(t *testing.T) {
	e := setupEcho()

	reqBody := model.CreateUserRequest{
		Name:     "Andi Wijaya",
		Username: "andi123",
		Password: "SecretPassword123!",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	eCtx := e.NewContext(req, rec)

	mockSvc := new(MockUserService)
	createdUser := model.User{
		ID:       10,
		Name:     reqBody.Name,
		Username: reqBody.Username,
	}

	mockSvc.On("CreateUser", mock.Anything, reqBody).Return(
		response.ResponseGenerate(http.StatusCreated, nil, nil, createdUser, nil),
		nil,
	).Once()

	ctrl := Controller{
		User:   mockSvc,
		Logger: zap.NewNop(),
	}

	err := ctrl.CreateUser(eCtx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var respBody response.FinalResponse
	err = json.Unmarshal(rec.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, 201, respBody.ResponseCode)

	mockSvc.AssertExpectations(t)
}

// Skenario 4: POST /v1/register Gagal Validasi (Bad Request 400)
// Ekspektasi: Payload tidak lengkap (misal password kosong) -> ditolak validator -> return 400 Bad Request.
// Service TIDAK boleh dipanggil.
func TestCreateUser_BadRequest_ValidationFailed(t *testing.T) {
	e := setupEcho()

	// Request tidak ada password
	incompleteBody := map[string]string{
		"name":     "Andi Wijaya",
		"username": "andi123",
		"password": "", // Kosong
	}
	bodyBytes, _ := json.Marshal(incompleteBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	eCtx := e.NewContext(req, rec)

	mockSvc := new(MockUserService)

	ctrl := Controller{
		User:   mockSvc,
		Logger: zap.NewNop(),
	}

	err := ctrl.CreateUser(eCtx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Pastikan Service TIDAK pernah dipanggil karena sudah gagal di validasi controller
	mockSvc.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

