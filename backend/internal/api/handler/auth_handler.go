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
	Email    string `json:"email" binding:"required,email"`
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

	accessToken, refreshToken, err := h.AuthService.Login(input.Username, input.Password)
	if err != nil {
		error.RespondError(c, error.NewUnauthorizedError(err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "bearer",
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input RefreshRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondError(c, error.NewBadRequestError("invalid request body"))
		return
	}

	accessToken, refreshToken, err := h.AuthService.RefreshToken(input.RefreshToken)
	if err != nil {
		error.RespondError(c, error.NewUnauthorizedError(err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "bearer",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.AuthService.Logout(userID); err != nil {
		error.RespondError(c, error.NewInternalError("failed to logout"))
		return
	}

	c.JSON(200, gin.H{"message": "logged out successfully"})
}
