package handler

import (
	"life-financial-assistant-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	BillService *service.BillService
}

type BillCreate struct {
	Amount     float64 `json:"amount" binding:"required"`
	Category   string  `json:"category" binding:"required"`
	Type       string  `json:"type" binding:"required"` // income/expense
	Remarks    string  `json:"description"`
	Date       string  `json:"bill_date"`
	IsPersonal bool    `json:"is_personal"`
	SpaceID    *uint   `json:"space_id"`
}

func (h *BillHandler) CreateBill(c *gin.Context) {
	var input BillCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	bill, err := h.BillService.CreateBill(userID, input.Amount, input.Category, input.Type, input.Remarks, input.IsPersonal, input.SpaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bill)
}

func (h *BillHandler) GetBills(c *gin.Context) {
	userID := c.GetUint("user_id")
	isPersonalStr := c.DefaultQuery("is_personal", "true")
	isPersonal, _ := strconv.ParseBool(isPersonalStr)
	spaceIDStr := c.Query("space_id")
	billType := c.Query("type")
	category := c.Query("category")

	var spaceID *uint
	if spaceIDStr != "" {
		id, err := strconv.ParseUint(spaceIDStr, 10, 32)
		if err == nil {
			val := uint(id)
			spaceID = &val
		}
	}

	bills, err := h.BillService.GetUserBills(userID, isPersonal, spaceID, billType, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bills)
}

func (h *BillHandler) GetBill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bill id"})
		return
	}

	userID := c.GetUint("user_id")

	bill, err := h.BillService.GetBillByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bill not found"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) UpdateBill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bill id"})
		return
	}

	var input BillCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	bill, err := h.BillService.UpdateBill(userID, uint(id), input.Amount, input.Category, input.Type, input.Remarks, input.IsPersonal, input.SpaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) GenerateAnalysis(c *gin.Context) {
	var input struct {
		IsPersonal bool  `json:"is_personal"`
		SpaceID    *uint `json:"space_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	report, err := h.BillService.GenerateAnalysis(userID, input.IsPersonal, input.SpaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

func (h *BillHandler) GetReports(c *gin.Context) {
	userID := c.GetUint("user_id")
	isPersonalStr := c.DefaultQuery("is_personal", "true")
	isPersonal, _ := strconv.ParseBool(isPersonalStr)
	spaceIDStr := c.Query("space_id")

	var spaceID *uint
	if spaceIDStr != "" {
		id, err := strconv.ParseUint(spaceIDStr, 10, 32)
		if err == nil {
			val := uint(id)
			spaceID = &val
		}
	}

	reports, err := h.BillService.GetReports(userID, isPersonal, spaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}
