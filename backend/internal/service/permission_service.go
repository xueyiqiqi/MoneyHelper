package service

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type PermissionService struct {
	Enforcer *casbin.Enforcer
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	adapter := gormadapter.NewAdapterByDB(db)

	e, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}

	// 加载策略
	err = e.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	// 初始化默认策略
	initPolicies(e)

	return &PermissionService{Enforcer: e}
}

func initPolicies(e *casbin.Enforcer) {
	policies := [][]string{
		// creator 可以对账单进行所有操作
		{"creator", "bill", "read"},
		{"creator", "bill", "write"},
		{"creator", "bill", "delete"},
		{"creator", "space", "admin"},

		// admin 可以对账单进行所有操作
		{"admin", "bill", "read"},
		{"admin", "bill", "write"},
		{"admin", "bill", "delete"},
		{"admin", "space", "admin"},

		// member 可以读写账单
		{"member", "bill", "read"},
		{"member", "bill", "write"},
		{"member", "space", "read"},

		// observer 只能读账单
		{"observer", "bill", "read"},
		{"observer", "space", "read"},
	}

	for _, p := range policies {
		e.AddPolicy(p[0], p[1], p[2])
	}

	e.SavePolicy()
}

// CheckPermission 检查角色是否有权限执行操作
func (s *PermissionService) CheckPermission(role, object, action string) bool {
	return s.Enforcer.HasPermission(role, object, action)
}
