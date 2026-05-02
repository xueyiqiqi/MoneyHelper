package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"life-financial-assistant-backend/internal/dto"
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

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.AuthService.GetProfile(userID)
	if err != nil {
		error.RespondError(c, error.NewNotFoundError("user not found"))
		return
	}

	c.JSON(200, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondError(c, error.NewBadRequestError("invalid request body"))
		return
	}

	user, err := h.AuthService.UpdateProfile(userID, input.Username, input.Email)
	if err != nil {
		error.RespondError(c, error.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(200, user)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondError(c, error.NewBadRequestError("invalid request body"))
		return
	}

	if err := h.AuthService.ChangePassword(userID, input.CurrentPassword, input.NewPassword); err != nil {
		error.RespondError(c, error.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(200, gin.H{"message": "password updated successfully"})
}

func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, err := c.FormFile("avatar")
	if err != nil {
		error.RespondError(c, error.NewBadRequestError("avatar file is required"))
		return
	}

	if file.Size > 2*1024*1024 {
		error.RespondError(c, error.NewBadRequestError("avatar file size must be 2MB or less"))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		error.RespondError(c, error.NewBadRequestError("unsupported avatar file type"))
		return
	}

	uploadDir := filepath.Join("uploads", "avatars")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		error.RespondError(c, error.NewInternalError("failed to prepare upload directory"))
		return
	}

	fileName := fmt.Sprintf("%d-%d%s", userID, time.Now().UnixNano(), ext)
	savedPath := filepath.Join(uploadDir, fileName)
	if err := c.SaveUploadedFile(file, savedPath); err != nil {
		error.RespondError(c, error.NewInternalError("failed to save avatar"))
		return
	}

	user, err := h.AuthService.UpdateAvatar(userID, "/uploads/avatars/"+fileName)
	if err != nil {
		error.RespondError(c, error.NewInternalError("failed to update avatar"))
		return
	}

	c.JSON(200, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.AuthService.Logout(userID); err != nil {
		error.RespondError(c, error.NewInternalError("failed to logout"))
		return
	}

	c.JSON(200, gin.H{"message": "logged out successfully"})
}
