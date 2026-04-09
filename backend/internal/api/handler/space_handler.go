package handler

import (
	"life-financial-assistant-backend/internal/service"
	"net/http"
	"strconv"

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

type AddMemberRequest struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (h *SpaceHandler) AddMember(c *gin.Context) {
	var input AddMemberRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spaceIDStr := c.Param("id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	if err := h.SpaceService.AddMember(uint(spaceID), input.Username, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}

func (h *SpaceHandler) LeaveSpace(c *gin.Context) {
	spaceIDStr := c.Param("id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	userID := c.GetUint("user_id")

	if err := h.SpaceService.LeaveSpace(uint(spaceID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left space successfully"})
}
