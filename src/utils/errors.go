package utils

import "fmt"

type AppError struct {
	Code       string
	Message    string
	HttpStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func ErrNotFound(msg string, args ...any) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: 404,
	}
}

func ErrConflict(msg string, args ...any) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: 409,
	}
}

func ErrValidation(msg string, args ...any) *AppError {
	return &AppError{
		Code:       "VALIDATION",
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: 422,
	}
}

func ErrBadRequest(msg string, args ...any) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: 400,
	}
}

func ErrInternal(msg string, args ...any) *AppError {
	return &AppError{
		Code:       "INTERNAL",
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: 500,
	}
}

func WrapInternal(err error, msg string) *AppError {
	return &AppError{
		Code:       "INTERNAL",
		Message:    msg,
		HttpStatus: 500,
		Err:        err,
	}
}
