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

type MockRoomService struct {
	CreateFunc      func(ctx context.Context, name string, capacitity int8) (models.Room, error)
	DeleteFunc      func(ctx context.Context, id int) (bool, error)
	EditCapacityFunc func(ctx context.Context, id int, capacitity int8) (models.Room, error)
	ListFunc        func(ctx context.Context) ([]models.Room, error)
}

func (m *MockRoomService) CreateRoomService(ctx context.Context, name string, capacitity int8) (models.Room, error) {
	return m.CreateFunc(ctx, name, capacitity)
}

func (m *MockRoomService) DeleteRoomService(ctx context.Context, id int) (bool, error) {
	return m.DeleteFunc(ctx, id)
}

func (m *MockRoomService) EditCapacitiyRoomService(ctx context.Context, id int, capacitity int8) (models.Room, error) {
	return m.EditCapacityFunc(ctx, id, capacitity)
}

func (m *MockRoomService) ListRoomsService(ctx context.Context) ([]models.Room, error) {
	return m.ListFunc(ctx)
}

func setupRoomTestContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/room", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func TestCreateRoom_Success(t *testing.T) {
	mock := &MockRoomService{
		CreateFunc: func(ctx context.Context, name string, capacitity int8) (models.Room, error) {
			return models.Room{ID: 1, Name: name, Capacitity: capacitity, IsActive: false}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupRoomTestContext(`{"name": "Sala 1", "capacitity": 10}`)

	controller.CreateRoom(c)

	assert.Equal(t, 201, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Sala 1", response["name"])
	assert.Equal(t, float64(10), response["capacitity"])
	assert.Equal(t, false, response["is_active"])
}

func TestCreateRoom_ValidationFailed(t *testing.T) {
	mock := &MockRoomService{
		CreateFunc: func(ctx context.Context, name string, capacitity int8) (models.Room, error) {
			t.Fatal("service não deveria ser chamado")
			return models.Room{}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupRoomTestContext(`{"name": "Sa", "capacitity": 10}`)

	controller.CreateRoom(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION", response.Code)
	assert.Contains(t, response.Error, "Dados invalidos")
}

func setupDeleteRoomContext(id string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("DELETE", "/room/delete/"+id, nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}

	return c, w
}

func setupListRoomsContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/room", nil)
	c.Request = req

	return c, w
}

func TestDeleteRoom_Success(t *testing.T) {
	mock := &MockRoomService{
		DeleteFunc: func(ctx context.Context, id int) (bool, error) {
			assert.Equal(t, 1, id)
			return true, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupDeleteRoomContext("1")

	controller.DeleteRoom(c)

	assert.Equal(t, 200, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Sala deletada com sucesso!", response["message"])
}

func TestDeleteRoom_InvalidID(t *testing.T) {
	mock := &MockRoomService{
		DeleteFunc: func(ctx context.Context, id int) (bool, error) {
			t.Fatal("service não deveria ser chamado")
			return false, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupDeleteRoomContext("abc")

	controller.DeleteRoom(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "BAD_REQUEST", response.Code)
	assert.Equal(t, "Id invalido", response.Error)
}

func setupEditCapacityContext(id string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PUT", "/room/"+id, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}

	return c, w
}

func TestEditCapacityRoom_Success(t *testing.T) {
	mock := &MockRoomService{
		EditCapacityFunc: func(ctx context.Context, id int, capacitity int8) (models.Room, error) {
			assert.Equal(t, 1, id)
			assert.Equal(t, int8(20), capacitity)
			return models.Room{ID: 1, Name: "Sala 1", Capacitity: capacitity, IsActive: true}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupEditCapacityContext("1", `{"capacitity": 20}`)

	controller.EditCapacitityRoom(c)

	assert.Equal(t, 200, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Sala 1", response["name"])
	assert.Equal(t, float64(20), response["capacitity"])
	assert.Equal(t, true, response["is_active"])
}

func TestEditCapacityRoom_ValidationFailed(t *testing.T) {
	mock := &MockRoomService{
		EditCapacityFunc: func(ctx context.Context, id int, capacitity int8) (models.Room, error) {
			t.Fatal("service não deveria ser chamado")
			return models.Room{}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupEditCapacityContext("1", `{"capacitity": 0}`)

	controller.EditCapacitityRoom(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION", response.Code)
	assert.Contains(t, response.Error, "Dados invalidos")
}

func TestEditCapacityRoom_InvalidID(t *testing.T) {
	mock := &MockRoomService{
		EditCapacityFunc: func(ctx context.Context, id int, capacitity int8) (models.Room, error) {
			t.Fatal("service não deveria ser chamado")
			return models.Room{}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupEditCapacityContext("abc", `{"capacitity": 10}`)

	controller.EditCapacitityRoom(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "BAD_REQUEST", response.Code)
	assert.Equal(t, "Id invalido", response.Error)
}

func TestListRooms_Success(t *testing.T) {
	mock := &MockRoomService{
		ListFunc: func(ctx context.Context) ([]models.Room, error) {
			return []models.Room{
				{ID: 1, Name: "Sala 1", Capacitity: 10, IsActive: true},
				{ID: 2, Name: "Sala 2", Capacitity: 5, IsActive: false},
			}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupListRoomsContext()

	controller.ListRooms(c)

	assert.Equal(t, 200, w.Code)

	var response []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, float64(1), response[0]["id"])
	assert.Equal(t, "Sala 1", response[0]["name"])
	assert.Equal(t, float64(10), response[0]["capacitity"])
	assert.Equal(t, true, response[0]["is_active"])
	assert.Equal(t, float64(2), response[1]["id"])
	assert.Equal(t, "Sala 2", response[1]["name"])
}

func TestListRooms_Empty(t *testing.T) {
	mock := &MockRoomService{
		ListFunc: func(ctx context.Context) ([]models.Room, error) {
			return []models.Room{}, nil
		},
	}

	controller := handlers.NewRoomController(mock)
	c, w := setupListRoomsContext()

	controller.ListRooms(c)

	assert.Equal(t, 200, w.Code)

	var response []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 0)
}
