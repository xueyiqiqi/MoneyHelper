package handler

import (
	"life-financial-assistant-backend/internal/error"
	"life-financial-assistant-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

type UserCreate struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input UserCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondError(c, error.NewBadRequestError("invalid request body"))
		return
	}

	if err := h.AuthService.Register(input.Username, input.Password, input.Email); err != nil {
		error.RespondError(c, error.NewInternalError("failed to register user"))
		return
	}

	c.JSON(201, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondError(c, error.NewBadRequestError("invalid request body"))
		return
	}

	token, err := h.AuthService.Login(input.Username, input.Password)
	if err != nil {
		error.RespondError(c, error.NewUnauthorizedError("invalid credentials"))
		return
	}

	c.JSON(200, gin.H{"access_token": token, "token_type": "bearer"})
}
