package error

import (
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func NewUnauthorizedError(msg string) *AppError {
	return &AppError{Code: 401, Message: msg}
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{Code: 404, Message: msg}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Code: 500, Message: msg}
}

func RespondError(c *gin.Context, err *AppError) {
	c.JSON(err.Code, gin.H{"error": err.Message, "code": err.Code})
}