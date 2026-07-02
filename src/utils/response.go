package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func JSONError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HttpStatus, ErrorResponse{
			Error: appErr.Message,
			Code:  appErr.Code,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "Erro interno do servidor",
		Code:  "INTERNAL",
	})
}

func JSONBindError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make([]string, 0)
		for _, fe := range validationErrors {
			fields = append(fields, fieldMessage(fe))
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Dados invalidos: " + strings.Join(fields, "; "),
			Code:  "VALIDATION",
		})
		return
	}
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: "Dados invalidos",
		Code:  "BAD_REQUEST",
	})
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " e obrigatorio"
	case "min":
		return fe.Field() + " deve ter no minimo " + fe.Param() + " caracteres"
	case "max":
		return fe.Field() + " deve ter no maximo " + fe.Param() + " caracteres"
	case "email":
		return fe.Field() + " deve ser um email valido"
	default:
		return fe.Field() + " invalido"
	}
}
