package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"teste-golang-api/src/handlers"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockBookingService struct {
	CreateFunc                func(ctx context.Context, roomID, userID int, title string, capacitityPeople int8, startDate, endDate time.Time) (models.Booking, error)
	ListFunc                  func(ctx context.Context, status string, page, limit int) (utils.PaginatedResponse, error)
	ListByUserFunc            func(ctx context.Context, userName string) ([]models.Booking, error)
	ListByRoomFunc            func(ctx context.Context, roomID int) ([]models.Booking, error)
	UpdateBookingStatusFunc   func(ctx context.Context, bookingID int) (bool, error)
}

func (m *MockBookingService) CreateBookingService(ctx context.Context, roomID, userID int, title string, capacitityPeople int8, startDate, endDate time.Time) (models.Booking, error) {
	return m.CreateFunc(ctx, roomID, userID, title, capacitityPeople, startDate, endDate)
}

func (m *MockBookingService) ListBookingsService(ctx context.Context, status string, page, limit int) (utils.PaginatedResponse, error) {
	return m.ListFunc(ctx, status, page, limit)
}

func (m *MockBookingService) ListBookingByUserService(ctx context.Context, userName string) ([]models.Booking, error) {
	return m.ListByUserFunc(ctx, userName)
}

func (m *MockBookingService) ListBookingsByRoomService(ctx context.Context, roomID int) ([]models.Booking, error) {
	return m.ListByRoomFunc(ctx, roomID)
}

func (m *MockBookingService) UpdateBookingStatusService(ctx context.Context, bookingID int) (bool, error) {
	return m.UpdateBookingStatusFunc(ctx, bookingID)
}

func setupBookingTestContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/booking", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func setupListBookingsContext(query string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/booking?"+query, nil)
	c.Request = req

	return c, w
}

func setupUpdateBookingContext(bookingID string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PUT", "/booking/"+bookingID, nil)
	c.Params = gin.Params{{Key: "id", Value: bookingID}}
	c.Request = req

	return c, w
}

func TestCreateBooking_Success(t *testing.T) {
	mock := &MockBookingService{
		CreateFunc: func(ctx context.Context, roomID, userID int, title string, capacitityPeople int8, startDate, endDate time.Time) (models.Booking, error) {
			assert.Equal(t, 1, roomID)
			assert.Equal(t, 1, userID)
			assert.Equal(t, "Reunião", title)
			assert.Equal(t, int8(5), capacitityPeople)
			return models.Booking{
				ID:               1,
				RoomID:           roomID,
				UserID:           userID,
				Title:            title,
				CapacitityPeople: capacitityPeople,
				StartDate:        startDate,
				EndDate:          endDate,
				Status:           true,
			}, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupBookingTestContext(`{
		"room_id": 1,
		"user_id": 1,
		"title": "Reunião",
		"capacitity_people": 5,
		"start_date": "2024-01-15",
		"end_date": "2024-01-16"
	}`)

	controller.CreateBooking(c)

	assert.Equal(t, 201, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, float64(1), response["room_id"])
	assert.Equal(t, float64(1), response["user_id"])
	assert.Equal(t, "Reunião", response["title"])
	assert.Equal(t, float64(5), response["capacitity_people"])
}

func TestCreateBooking_ValidationFailed(t *testing.T) {
	mock := &MockBookingService{
		CreateFunc: func(ctx context.Context, roomID, userID int, title string, capacitityPeople int8, startDate, endDate time.Time) (models.Booking, error) {
			t.Fatal("service não deveria ser chamado")
			return models.Booking{}, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupBookingTestContext(`{
		"room_id": 1,
		"user_id": 1,
		"title": "Re",
		"capacitity_people": 5,
		"start_date": "2024-01-15",
		"end_date": "2024-01-16"
	}`)

	controller.CreateBooking(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION", response.Code)
	assert.Contains(t, response.Error, "Dados invalidos")
}

func TestCreateBooking_InvalidDateFormat(t *testing.T) {
	mock := &MockBookingService{
		CreateFunc: func(ctx context.Context, roomID, userID int, title string, capacitityPeople int8, startDate, endDate time.Time) (models.Booking, error) {
			t.Fatal("service não deveria ser chamado")
			return models.Booking{}, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupBookingTestContext(`{
		"room_id": 1,
		"user_id": 1,
		"title": "Reunião",
		"capacitity_people": 5,
		"start_date": "15/01/2024",
		"end_date": "2024-01-16"
	}`)

	controller.CreateBooking(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "BAD_REQUEST", response.Code)
	assert.Contains(t, response.Error, "start_date invalido")
}

func TestListBookings_Success(t *testing.T) {
	startDate, _ := time.Parse("2006-01-02", "2024-01-15")
	endDate, _ := time.Parse("2006-01-02", "2024-01-16")

	mock := &MockBookingService{
		ListFunc: func(ctx context.Context, status string, page, limit int) (utils.PaginatedResponse, error) {
			assert.Equal(t, "", status)
			assert.Equal(t, 1, page)
			assert.Equal(t, 10, limit)
			return utils.PaginatedResponse{
				Data: []models.Booking{
					{ID: 1, RoomID: 1, UserID: 1, Title: "Reunião 1", CapacitityPeople: 5, StartDate: startDate, EndDate: endDate, Status: true},
					{ID: 2, RoomID: 2, UserID: 1, Title: "Reunião 2", CapacitityPeople: 3, StartDate: startDate, EndDate: endDate, Status: true},
				},
				Page:       1,
				Limit:      10,
				Total:      2,
				TotalPages: 1,
			}, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupListBookingsContext("")

	controller.ListBookings(c)

	assert.Equal(t, 200, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(10), response["limit"])
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["total_pages"])

	data, ok := response["data"].([]any)
	require.True(t, ok)
	assert.Len(t, data, 2)
}

func TestListBookings_WithStatusFilter(t *testing.T) {
	mock := &MockBookingService{
		ListFunc: func(ctx context.Context, status string, page, limit int) (utils.PaginatedResponse, error) {
			assert.Equal(t, "active", status)
			assert.Equal(t, 2, page)
			assert.Equal(t, 5, limit)
			return utils.PaginatedResponse{
				Data:       []models.Booking{},
				Page:       2,
				Limit:      5,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupListBookingsContext("status=active&page=2&limit=5")

	controller.ListBookings(c)

	assert.Equal(t, 200, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(2), response["page"])
	assert.Equal(t, float64(5), response["limit"])
	assert.Equal(t, float64(0), response["total"])
}

func TestUpdateBookingStatus_Success(t *testing.T) {
	mock := &MockBookingService{
		UpdateBookingStatusFunc: func(ctx context.Context, bookingID int) (bool, error) {
			assert.Equal(t, 1, bookingID)
			return true, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupUpdateBookingContext("1")

	controller.UpdateBookingStatus(c)

	assert.Equal(t, 200, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, true, response["success"])
}

func TestUpdateBookingStatus_NotFound(t *testing.T) {
	mock := &MockBookingService{
		UpdateBookingStatusFunc: func(ctx context.Context, bookingID int) (bool, error) {
			return false, utils.ErrNotFound("Reserva nao existe ou ja esta inativa")
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupUpdateBookingContext("999")

	controller.UpdateBookingStatus(c)

	assert.Equal(t, 404, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "NOT_FOUND", response.Code)
	assert.Contains(t, response.Error, "Reserva nao existe")
}

func TestUpdateBookingStatus_InvalidID(t *testing.T) {
	mock := &MockBookingService{
		UpdateBookingStatusFunc: func(ctx context.Context, bookingID int) (bool, error) {
			t.Fatal("service não deveria ser chamado")
			return false, nil
		},
	}

	controller := handlers.NewBookingController(mock)
	c, w := setupUpdateBookingContext("abc")

	controller.UpdateBookingStatus(c)

	assert.Equal(t, 400, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "BAD_REQUEST", response.Code)
	assert.Contains(t, response.Error, "id invalido")
}
