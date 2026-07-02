package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"teste-golang-api/src/handlers"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockUserService struct {
	CreateFunc func(ctx context.Context, name, email string) (models.User, error)
	DeleteFunc func(ctx context.Context, id int) (bool, error)
	ListFunc   func(ctx context.Context) ([]models.User, error)
}

func (m *MockUserService) CreateUserService(ctx context.Context, name, email string) (models.User, error) {
	return m.CreateFunc(ctx, name, email)
}

func (m *MockUserService) DeleteUserService(ctx context.Context, id int) (bool, error) {
	return m.DeleteFunc(ctx, id)
}

func (m *MockUserService) ListUsersService(ctx context.Context) ([]models.User, error) {
	return m.ListFunc(ctx)
}

func setupTestContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func TestCreateUser_Success(t *testing.T) {
	mock := &MockUserService{
		CreateFunc: func(ctx context.Context, name, email string) (models.User, error) {
			return models.User{ID: 1, Name: name}, nil
		},
	}

	controller := handlers.NewUserController(mock)
	c, w := setupTestContext(`{"name": "Lucas", "email": "lucas@test.com"}`)

	controller.CreateUser(c)

	assert.Equal(t, 201, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Lucas", response["name"])
}

func TestCreateUser_ValidationFailed(t *testing.T) {
	mock := &MockUserService{
		CreateFunc: func(ctx context.Context, name, email string) (models.User, error) {
			t.Fatal("service não deveria ser chamado")
			return models.User{}, nil
		},
	}

	controller := handlers.NewUserController(mock)
	c, w := setupTestContext(`{"name": "Lu", "email": "lucas@test.com"}`)

	controller.CreateUser(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION", response.Code)
	assert.Contains(t, response.Error, "Dados invalidos")
}

func setupDeleteContext(id string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("DELETE", "/user/delete/"+id, nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}

	return c, w
}

func setupListUsersContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/user", nil)
	c.Request = req

	return c, w
}

func TestDeleteUser_Success(t *testing.T) {
	mock := &MockUserService{
		DeleteFunc: func(ctx context.Context, id int) (bool, error) {
			assert.Equal(t, 1, id)
			return true, nil
		},
	}

	controller := handlers.NewUserController(mock)
	c, w := setupDeleteContext("1")

	controller.DeleteUser(c)

	assert.Equal(t, 200, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Usuário deletado com sucesso!", response["message"])
}

func TestDeleteUser_InvalidID(t *testing.T) {
	mock := &MockUserService{
		DeleteFunc: func(ctx context.Context, id int) (bool, error) {
			t.Fatal("service não deveria ser chamado")
			return false, nil
		},
	}

	controller := handlers.NewUserController(mock)
	c, w := setupDeleteContext("abc")

	controller.DeleteUser(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "BAD_REQUEST", response.Code)
	assert.Equal(t, "Id inválido", response.Error)
}

func TestListUsers_Success(t *testing.T) {
	mock := &MockUserService{
		ListFunc: func(ctx context.Context) ([]models.User, error) {
			return []models.User{
				{ID: 1, Name: "Lucas", Email: "lucas@test.com"},
				{ID: 2, Name: "Maria", Email: "maria@test.com"},
			}, nil
		},
	}

	controller := handlers.NewUserController(mock)
	c, w := setupListUsersContext()

	controller.ListUsers(c)

	assert.Equal(t, 200, w.Code)

	var response []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, float64(1), response[0]["id"])
	assert.Equal(t, "Lucas", response[0]["name"])
	assert.Equal(t, "lucas@test.com", response[0]["email"])
	assert.Equal(t, float64(2), response[1]["id"])
	assert.Equal(t, "Maria", response[1]["name"])
}

func TestListUsers_Empty(t *testing.T) {
	mock := &MockUserService{
		ListFunc: func(ctx context.Context) ([]models.User, error) {
			return []models.User{}, nil
		},
	}

	controller := handlers.NewUserController(mock)
	c, w := setupListUsersContext()

	controller.ListUsers(c)

	assert.Equal(t, 200, w.Code)

	var response []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 0)
}
