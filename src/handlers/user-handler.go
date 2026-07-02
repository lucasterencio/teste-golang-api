package handlers

import (
	"context"
	"strconv"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUserService(ctx context.Context, name string, email string) (models.User, error)
	DeleteUserService(ctx context.Context, user_id int) (bool, error)
	ListUsersService(ctx context.Context) ([]models.User, error)
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=3,max=100"`
	Email string `json:"email" binding:"required,email"`
}

type UserController struct {
	service UserService
}

func NewUserController(s UserService) *UserController {
	return &UserController{service: s}
}

// @Summary      Criar Usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      CreateUserRequest true "Dados do usuário"
// @Success      201   {object}  utils.UserResponse "Usuário criado com sucesso"
// @Failure      400   {object}  utils.ErrorResponse "Dados inválidos"
// @Failure      500   {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /user [post]
func (controller *UserController) CreateUser(c *gin.Context) {
	var body CreateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.JSONBindError(c, err)
		return
	}

	user, err := controller.service.CreateUserService(c.Request.Context(), body.Name, body.Email)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 201, gin.H{
		"id":   user.ID,
		"name": user.Name,
	})
}

// @Summary      Deletar Usuário
// @Description  Remove um usuário do sistema pelo ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "ID do usuário"
// @Success      200  {object}  utils.SuccessResponse "Usuário deletado com sucesso"
// @Failure      400  {object}  utils.ErrorResponse "ID inválido"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /user/delete/{id} [delete]
func (controller *UserController) DeleteUser(c *gin.Context){
	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userIDStr == "" {
		utils.JSONError(c, utils.ErrBadRequest("Id inválido"))
		return
	}

	_, err = controller.service.DeleteUserService(c.Request.Context(), userID)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, gin.H{
		"message": "Usuário deletado com sucesso!",
	})
}

// ListUsersHandler lista todos os usuários cadastrados
// @Summary      Listar Usuários
// @Description  Retorna uma lista de todos os usuários cadastrados no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User "Lista de usuários retornada com sucesso"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /user [get]
func (controller *UserController) ListUsers(c *gin.Context) {
	users, err := controller.service.ListUsersService(c)

	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSON(c, 200, users)
}