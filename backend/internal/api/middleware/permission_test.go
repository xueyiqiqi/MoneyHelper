package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appmodel "life-financial-assistant-backend/internal/model"
	"life-financial-assistant-backend/internal/repository"
	"life-financial-assistant-backend/internal/service"

	"github.com/casbin/casbin/v3"
	casbinmodel "github.com/casbin/casbin/v3/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupPermissionTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&appmodel.User{}, &appmodel.Space{}, &appmodel.SpaceUserLink{}, &appmodel.Bill{}); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	return db
}

func setupPermissionTestService(t *testing.T) *service.PermissionService {
	t.Helper()

	m, err := casbinmodel.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		t.Fatalf("build casbin model: %v", err)
	}

	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("create enforcer: %v", err)
	}

	_, _ = e.AddPolicy("member", "bill", "read")
	_, _ = e.AddPolicy("member", "bill", "write")
	_, _ = e.AddPolicy("member", "space", "read")
	_, _ = e.AddPolicy("observer", "bill", "read")
	_, _ = e.AddPolicy("observer", "space", "read")
	_, _ = e.AddPolicy("admin", "space", "admin")

	return &service.PermissionService{Enforcer: e}
}

func TestPermissionMiddlewareRejectsNonMemberForSpaceBillRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	r := gin.New()
	r.GET("/bills", func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	}, PermissionMiddleware("bill", "read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/bills?space_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-member space bill read, got %d", w.Code)
	}
}

func TestBillPermissionMiddlewareRejectsUnauthorizedPersonalBillRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	repository.DB = db
	defer func() {
		repository.DB = originalDB
	}()

	owner := appmodel.User{Username: "owner", Email: "owner@example.com", HashedPassword: "x"}
	viewer := appmodel.User{Username: "viewer", Email: "viewer@example.com", HashedPassword: "x"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("create owner: %v", err)
	}
	if err := db.Create(&viewer).Error; err != nil {
		t.Fatalf("create viewer: %v", err)
	}

	bill := appmodel.Bill{
		Amount:     10,
		Category:   "food",
		Type:       "income",
		Remarks:    "test",
		Date:       time.Now(),
		IsPersonal: true,
		UserID:     owner.ID,
	}
	if err := db.Create(&bill).Error; err != nil {
		t.Fatalf("create bill: %v", err)
	}

	r := gin.New()
	r.GET("/bills/:id", func(c *gin.Context) {
		c.Set("user_id", viewer.ID)
		c.Next()
	}, BillPermissionMiddleware("read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/bills/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for unauthorized personal bill read, got %d", w.Code)
	}
}


func TestPermissionMiddlewareAllowsMemberToReadSpaceBills(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	user := appmodel.User{Username: "member-reader", Email: "member-reader@example.com", HashedPassword: "x"}
	space := appmodel.Space{Name: "space-read"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&appmodel.SpaceUserLink{UserID: user.ID, SpaceID: space.ID, Role: appmodel.RoleMember}).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}

	r := gin.New()
	r.GET("/bills", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Next()
	}, PermissionMiddleware("bill", "read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/bills?space_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for member space bill read, got %d", w.Code)
	}
}

func TestPermissionMiddlewareAllowsObserverToReadSpaceBills(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	user := appmodel.User{Username: "observer-reader", Email: "observer-reader@example.com", HashedPassword: "x"}
	space := appmodel.Space{Name: "space-observer-read"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&appmodel.SpaceUserLink{UserID: user.ID, SpaceID: space.ID, Role: appmodel.RoleObserver}).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}

	r := gin.New()
	r.GET("/bills", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Next()
	}, PermissionMiddleware("bill", "read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/bills?space_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for observer space bill read, got %d", w.Code)
	}
}

func TestPermissionMiddlewareAllowsMemberToWriteSpaceBills(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	user := appmodel.User{Username: "member-writer", Email: "member-writer@example.com", HashedPassword: "x"}
	space := appmodel.Space{Name: "space-write"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&appmodel.SpaceUserLink{UserID: user.ID, SpaceID: space.ID, Role: appmodel.RoleMember}).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}

	r := gin.New()
	r.POST("/bills", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Next()
	}, PermissionMiddleware("bill", "write"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/bills?space_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for member space bill write, got %d", w.Code)
	}
}

func TestPermissionMiddlewareRejectsObserverForSpaceBillWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	user := appmodel.User{Username: "observer-writer", Email: "observer-writer@example.com", HashedPassword: "x"}
	space := appmodel.Space{Name: "space-observer-write"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&appmodel.SpaceUserLink{UserID: user.ID, SpaceID: space.ID, Role: appmodel.RoleObserver}).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}

	r := gin.New()
	r.POST("/bills", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Next()
	}, PermissionMiddleware("bill", "write"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/bills?space_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for observer space bill write, got %d", w.Code)
	}
}


func TestPermissionMiddlewareAllowsMemberToReadSpaceMembers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	user := appmodel.User{Username: "member-list-reader", Email: "member-list-reader@example.com", HashedPassword: "x"}
	space := appmodel.Space{Name: "space-members-read"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&space).Error; err != nil {
		t.Fatalf("create space: %v", err)
	}
	if err := db.Create(&appmodel.SpaceUserLink{UserID: user.ID, SpaceID: space.ID, Role: appmodel.RoleMember}).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}

	r := gin.New()
	r.GET("/spaces/:id/members", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		c.Next()
	}, PermissionMiddleware("space", "read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/spaces/1/members", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for member reading member list, got %d", w.Code)
	}
}

func TestPermissionMiddlewareRejectsNonMemberForSpaceMembersRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupPermissionTestDB(t)
	originalDB := repository.DB
	originalPermissionService := permissionService
	repository.DB = db
	permissionService = setupPermissionTestService(t)
	defer func() {
		repository.DB = originalDB
		permissionService = originalPermissionService
	}()

	r := gin.New()
	r.GET("/spaces/:id/members", func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	}, PermissionMiddleware("space", "read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/spaces/1/members", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-member reading member list, got %d", w.Code)
	}
}
