package middleware

import (
	"net/http"
	"strconv"

	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
	"life-financial-assistant-backend/internal/service"

	"github.com/gin-gonic/gin"
)

var permissionService *service.PermissionService

func SetPermissionService(ps *service.PermissionService) {
	permissionService = ps
}

// PermissionMiddleware 创建权限检查中间件
func PermissionMiddleware(object, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		spaceIDStr := c.Param("id")
		if spaceIDStr == "" {
			spaceIDStr = c.Query("space_id")
		}

		var spaceID uint
		if spaceIDStr != "" {
			id, err := strconv.ParseUint(spaceIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
				c.Abort()
				return
			}
			spaceID = uint(id)
		}

		var role model.Role
		if spaceID > 0 {
			var link model.SpaceUserLink
			err := repository.DB.
				Where("user_id = ? AND space_id = ?", userID, spaceID).
				First(&link).Error
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
				c.Abort()
				return
			}
			role = link.Role
		} else {
			role = model.RoleMember
		}

		if !permissionService.CheckPermission(string(role), object, action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func BillPermissionMiddleware(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		billIDStr := c.Param("id")
		billID, err := strconv.ParseUint(billIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bill id"})
			c.Abort()
			return
		}

		bill, err := repository.BillAccessRepository{}.GetByID(uint(billID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bill not found"})
			c.Abort()
			return
		}

		if bill.IsPersonal {
			if bill.UserID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		if bill.SpaceID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		var link model.SpaceUserLink
		err = repository.DB.
			Where("user_id = ? AND space_id = ?", userID, *bill.SpaceID).
			First(&link).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		if !permissionService.CheckPermission(string(link.Role), "bill", action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
