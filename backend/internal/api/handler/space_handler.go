package handler

import (
	"life-financial-assistant-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	SpaceService *service.SpaceService
}

type SpaceCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	var input SpaceCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	space, err := h.SpaceService.CreateSpace(userID, input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, space)
}

func (h *SpaceHandler) GetUserSpaces(c *gin.Context) {
	userID := c.GetUint("user_id")

	spaces, err := h.SpaceService.GetUserSpaces(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, spaces)
}
