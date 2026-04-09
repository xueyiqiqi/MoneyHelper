package middleware

import (
	"net/http"
	"strconv"

	"life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

var PermissionService interface {
	CheckPermission(role, object, action string) bool
}

func SetPermissionService(ps interface{}) {
	PermissionService = ps
}

// PermissionMiddleware 创建权限检查中间件
func PermissionMiddleware(object, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		// 从请求参数或上下文获取 space_id
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

		// 获取用户在空间的角色
		var role model.Role
		if spaceID > 0 {
			var link model.SpaceUserLink
			err := repository.DB.
				Where("user_id = ? AND space_id = ?", userID, spaceID).
				First(&link).Error
			if err != nil {
				// 个人账单不需要 space_id，默认为 member
				role = model.RoleMember
			} else {
				role = link.Role
			}
		} else {
			// 没有 space_id 时，默认给 member 权限
			role = model.RoleMember
		}

		// 检查权限
		if !PermissionService.CheckPermission(string(role), object, action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}