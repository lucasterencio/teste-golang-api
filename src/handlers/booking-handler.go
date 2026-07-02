package handlers

import (
	"context"
	"strconv"
	"time"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/gin-gonic/gin"
)

type BookingService interface {
	CreateBookingService(ctx context.Context, room_id int, user_id int, title string, capacitity_people int8, startDate time.Time, endDate time.Time) (models.Booking, error)
	ListBookingsService(ctx context.Context, status string, page int, limit int) (utils.PaginatedResponse, error)
	ListBookingByUserService(ctx context.Context, user_name string) ([]models.Booking, error)
	ListBookingsByRoomService(ctx context.Context, room_id int) ([]models.Booking, error)
	UpdateBookingStatusService(ctx context.Context, booking_id int) (bool, error)
}

type CreateBookingRequest struct {
	RoomID           int    `json:"room_id" binding:"required"`
	UserID           int    `json:"user_id" binding:"required"`
	Title            string `json:"title" binding:"required,min=3,max=100"`
	CapacitityPeople int8   `json:"capacitity_people" binding:"required"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
}

type BookingController struct {
	service BookingService
}

func NewBookingController(s BookingService) *BookingController {
	return &BookingController{service: s}
}

// CreateBookingHandler cria uma nova reserva
// @Summary      Criar Reserva
// @Description  Cria uma nova reserva de sala
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        body  body      CreateBookingRequest true "Dados da reserva"
// @Success      201   {object}  models.Booking "Reserva criada com sucesso"
// @Failure      400   {object}  utils.ErrorResponse "Dados inválidos"
// @Failure      500   {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /booking [post]
func (controller *BookingController) CreateBooking(c *gin.Context) {
	var body CreateBookingRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONBindError(c, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		utils.JSONError(c, utils.ErrBadRequest("start_date invalido, use o formato YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		utils.JSONError(c, utils.ErrBadRequest("end_date invalido, use o formato YYYY-MM-DD"))
		return
	}

	booking, err := controller.service.CreateBookingService(
		c, body.RoomID, body.UserID, body.Title, body.CapacitityPeople, startDate, endDate,
	)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 201, booking)
}


// ListUsersHandler lista todos os usuários cadastrados
// @Summary      Listar Reservas
// @Description  Retorna uma lista de todas as reservas cadastradas no sistema
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Booking "Lista de reservas retornada com sucesso"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /booking [get]
func (controller *BookingController) ListBookings(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := controller.service.ListBookingsService(c, status, page, limit)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, result)
}

// ListBookingsByUserHandler lista reservas por usuário
// @Summary      Listar Reservas por Usuário
// @Description  Retorna todas as reservas de um usuário específico
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        username  query     string  true  "Nome do usuário"
// @Success      200       {array}   models.Booking "Lista de reservas do usuário"
// @Failure      500       {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /booking/by-user [get]
func (controller *BookingController) ListBookingsByUser(c *gin.Context) {
	userName := c.Query("username")

	bookings, err := controller.service.ListBookingByUserService(c, userName)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, bookings)
}

// ListBookingsByRoomHandler lista reservas por sala
// @Summary      Listar Reservas por Sala
// @Description  Retorna todas as reservas de uma sala específica
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        room_id  query     int  true  "ID da sala"
// @Success      200      {array}   models.Booking "Lista de reservas da sala"
// @Failure      400      {object}  utils.ErrorResponse "room_id inválido"
// @Failure      500      {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /booking/by-room [get]
func (controller *BookingController) ListBookingsByRoom(c *gin.Context) {
	roomIDStr := c.Query("room_id")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil || roomIDStr == "" {
		utils.JSONError(c, utils.ErrBadRequest("room_id invalido"))
		return
	}

	bookings, err := controller.service.ListBookingsByRoomService(c, roomID)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, bookings)
}

// UpdateBookingStatusHandler desativa uma reserva
// @Summary      Desativar Reserva
// @Description  Desativa uma reserva (soft delete), liberando a sala para novas reservas
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "ID da reserva"
// @Success      200  {object}  map[string]bool "Reserva desativada com sucesso"
// @Failure      400  {object}  utils.ErrorResponse "ID inválido"
// @Failure      404  {object}  utils.ErrorResponse "Reserva não encontrada"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /booking/{id} [put]
func (controller *BookingController) UpdateBookingStatus(c *gin.Context) {
	bookingIDStr := c.Param("id")

	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil || bookingIDStr == "" {
		utils.JSONError(c, utils.ErrBadRequest("id invalido"))
		return
	}

	result, err := controller.service.UpdateBookingStatusService(c, bookingID)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, map[string]bool{"success": result})
}
