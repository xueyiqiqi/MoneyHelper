package handler

import (
	"life-financial-assistant-backend/internal/error"
	"life-financial-assistant-backend/internal/model"
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
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
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

	if err := h.SpaceService.AddMember(uint(spaceID), input.Email, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (h *SpaceHandler) UpdateMemberRole(c *gin.Context) {
	var input UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spaceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.SpaceService.UpdateMemberRole(operatorID, uint(spaceID), uint(targetUserID), model.Role(input.Role)); err != nil {
		if appErr, ok := err.(*error.AppError); ok {
			error.RespondError(c, appErr)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member role updated successfully"})
}

func (h *SpaceHandler) RemoveMember(c *gin.Context) {
	spaceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.SpaceService.RemoveMember(operatorID, uint(spaceID), uint(targetUserID)); err != nil {
		if appErr, ok := err.(*error.AppError); ok {
			error.RespondError(c, appErr)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
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

func (h *SpaceHandler) GetMembers(c *gin.Context) {
	spaceIDStr := c.Param("id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	members, err := h.SpaceService.GetMembers(uint(spaceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}
