package service

import (
	"log"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type PermissionService struct {
	Enforcer *casbin.Enforcer
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}

	// 创建 enforcer (使用 v3，完美兼容 gorm-adapter/v3)
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
		// creator 可以对空间进行管理和读取
		{"creator", "space", "admin"},
		{"creator", "space", "read"},

		// admin 可以对账单进行所有操作
		{"admin", "bill", "read"},
		{"admin", "bill", "write"},
		{"admin", "bill", "delete"},
		{"admin", "space", "admin"},
		{"admin", "space", "read"},

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
	ok, err := s.Enforcer.Enforce(role, object, action)
	if err != nil {
		log.Printf("Enforce error: %v", err)
		return false
	}
	return ok
}
