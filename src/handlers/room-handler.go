package handlers

import (
	"context"
	"strconv"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/gin-gonic/gin"
)

type RoomService interface {
	CreateRoomService(ctx context.Context, name string, capacitity int8) (models.Room, error)
	DeleteRoomService(ctx context.Context, room_id int) (bool, error)
	EditCapacitiyRoomService(ctx context.Context, room_id int, capacitity int8) (models.Room, error)
	ListRoomsService(ctx context.Context) ([]models.Room, error)
}

type CreateRoomRequest struct {
	Name       string `json:"name" binding:"required,min=3,max=100"`
	Capacitity int8   `json:"capacitity" binding:"required,min=1"`
}

type RoomController struct {
	service RoomService
}

func NewRoomController(s RoomService) *RoomController {
	return &RoomController{service: s}
}

// @Summary      Criar Sala
// @Description  Cria uma nova sala no sistema
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRoomRequest true "Dados da sala"
// @Success      201   {object}  utils.SuccessResponse "Sala criada com sucesso"
// @Failure      400   {object}  utils.ErrorResponse "Dados inválidos"
// @Failure      500   {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /room [post]
func (controller *RoomController) CreateRoom(c *gin.Context) {
	var body CreateRoomRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONBindError(c, err)
		return
	}

	room, err := controller.service.CreateRoomService(c.Request.Context(), body.Name, body.Capacitity)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 201, gin.H{
		"id":         room.ID,
		"name":       room.Name,
		"is_active":  room.IsActive,
		"capacitity": room.Capacitity,
	})
}

// @Summary      Deletar Sala
// @Description  Remove uma sala do sistema pelo ID
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "ID da sala"
// @Success      200  {object}  utils.SuccessResponse "Sala deletada com sucesso"
// @Failure      400  {object}  utils.ErrorResponse "ID inválido"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /room/delete/{id} [delete]
func (controller *RoomController) DeleteRoom(c *gin.Context) {
	roomIDStr := c.Param("id")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil || roomIDStr == "" {
		utils.JSONError(c, utils.ErrBadRequest("Id invalido"))
		return
	}

	_, err = controller.service.DeleteRoomService(c.Request.Context(), roomID)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, gin.H{
		"message": "Sala deletada com sucesso!",
	})
}

type EditCapacitityRequest struct {
	Capacitity int8 `json:"capacitity" binding:"required,min=1"`
}

// @Summary      Editar Capacidade da Sala
// @Description  Atualiza a capacidade de uma sala existente
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "ID da sala"
// @Param        body  body      EditCapacitityRequest  true  "Nova capacidade"
// @Success      200   {object}  utils.SuccessResponse "Sala atualizada com sucesso"
// @Failure      400   {object}  utils.ErrorResponse "Dados inválidos"
// @Failure      500   {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /room/{id} [put]
func (controller *RoomController) EditCapacitityRoom(c *gin.Context) {
	var body EditCapacitityRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONBindError(c, err)
		return
	}

	roomIDStr := c.Param("id")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil || roomIDStr == "" {
		utils.JSONError(c, utils.ErrBadRequest("Id invalido"))
		return
	}

	room, err := controller.service.EditCapacitiyRoomService(c.Request.Context(), roomID, body.Capacitity)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, gin.H{
		"id":         room.ID,
		"name":       room.Name,
		"is_active":  room.IsActive,
		"capacitity": room.Capacitity,
	})
}

// ListRoomsHandler lista todas as salas cadastradas
// @Summary      Listar Salas
// @Description  Retorna uma lista de todas as salas cadastradas no sistema
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Room "Lista de salas retornada com sucesso"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /room [get]
func (controller *RoomController) ListRooms(c *gin.Context) {
	rooms, err := controller.service.ListRoomsService(c)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, rooms)
}