# 后端优化实施计划

**Goal:** 对 Go 后端进行架构优化，包括错误处理统一、中间件整理、类型完善、日志和环境配置

**Architecture:** 渐进式优化，从配置和错误处理开始，逐步改进代码结构

**Tech Stack:** Go + Gin + Gorm

---

## Task 1: 环境变量配置

**Files:**
- Create: `backend/internal/config/config.go`

- [ ] **Step 1: 创建配置模块**

```go
package config

import (
    "os"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    AI       AIConfig
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    Path string
}

type AIConfig struct {
    APIKey string
}

func Load() *Config {
    return &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
        Database: DatabaseConfig{
            Path: getEnv("DATABASE_PATH", "./database.db"),
        },
        AI: AIConfig{
            APIKey: getEnv("AI_API_KEY", ""),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

- [ ] **Step 2: 更新 main.go 使用配置**

```go
func main() {
    cfg := config.Load()
    
    repository.InitDB(cfg.Database.Path)
    // ... 传递 config 给需要的地方
}
```

- [ ] **Step 3: 提交**

```bash
git add backend/internal/config/config.go backend/cmd/main.go
git commit -m "feat: 添加环境变量配置模块"
```

---

## Task 2: 错误处理统一

**Files:**
- Create: `backend/internal/error/errors.go`
- Modify: `backend/internal/api/handler/auth_handler.go`

- [ ] **Step 1: 创建错误类型定义**

```go
package error

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewBadRequestError(msg string) *AppError {
    return &AppError{Code: 400, Message: msg}
}

func NewUnauthorizedError(msg string) *AppError {
    return &AppError{Code: 401, Message: msg}
}

func NewNotFoundError(msg string) *AppError {
    return &AppError{Code: 404, Message: msg}
}

func NewInternalError(msg string) *AppError {
    return &AppError{Code: 500, Message: msg}
}
```

- [ ] **Step 2: 创建响应助手**

```go
package error

import (
    "github.com/gin-gonic/gin"
)

func RespondError(c *gin.Context, err *AppError) {
    c.JSON(err.Code, gin.H{"error": err.Message, "code": err.Code})
}
```

- [ ] **Step 3: 在 handler 中使用**

修改 auth_handler.go:

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        error.RespondError(c, error.NewBadRequestError("invalid request body"))
        return
    }
    // ...
}
```

- [ ] **Step 4: 提交**

```bash
git add backend/internal/error/errors.go backend/internal/api/handler/auth_handler.go
git commit -m "feat: 统一错误处理类型"
```

---

## Task 3: 中间件整理

**Files:**
- Rename: `backend/internal/api/middleware/auth_middleware.go` → `backend/internal/api/middleware/auth.go`
- Rename: `backend/internal/api/middleware/package middleware.go` → `backend/internal/api/middleware/cors.go`
- Create: `backend/internal/api/middleware/recovery.go`

- [ ] **Step 1: 重命名文件**

```bash
cd backend/internal/api/middleware
mv "auth_middleware.go" "auth.go"
mv "package middleware.go" "cors.go"
```

- [ ] **Step 2: 添加 Recovery 中间件**

```go
package middleware

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "internal server error",
                    "code":  500,
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

- [ ] **Step 3: 更新 main.go 引入 recovery**

```go
r := gin.New()
r.Use(middleware.RecoveryMiddleware())
r.Use(middleware.CORSMiddleware())
```

- [ ] **Step 4: 提交**

```bash
git add backend/internal/api/middleware/
git commit -m "refactor: 规范化中间件文件命名"
```

---

## Task 4: 类型/模型完善

**Files:**
- Create: `backend/internal/dto/auth.go`
- Create: `backend/internal/dto/space.go`
- Create: `backend/internal/dto/bill.go`

- [ ] **Step 1: 创建 Auth DTO**

```go
package dto

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
}
```

- [ ] **Step 2: 创建 Space DTO**

```go
package dto

type CreateSpaceRequest struct {
    Name        string `json:"name" binding:"required,min=1,max=100"`
    Description string `json:"description"`
}

type SpaceResponse struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}
```

- [ ] **Step 3: 创建 Bill DTO**

```go
package dto

type CreateBillRequest struct {
    Amount     float64 `json:"amount" binding:"required,gt=0"`
    Category   string  `json:"category" binding:"required"`
    Remarks    string  `json:"remarks"`
    Date       string  `json:"date"`
    IsPersonal bool    `json:"is_personal"`
    SpaceID    *uint   `json:"space_id,omitempty"`
}

type BillResponse struct {
    ID          uint    `json:"id"`
    Amount     float64 `json:"amount"`
    Category   string  `json:"category"`
    Remarks    string  `json:"remarks,omitempty"`
    Date       string  `json:"date"`
    CreatorName string `json:"creator_name,omitempty"`
}
```

- [ ] **Step 4: 更新 handler 使用 DTO**

修改 handler 使用 dto 类型：

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        error.RespondError(c, error.NewBadRequestError("invalid request"))
        return
    }
    // use req.Username, req.Password
}
```

- [ ] **Step 5: 提交**

```bash
git add backend/internal/dto/
git commit -m "feat: 添加 DTO 数据传输对象"
```

---

## Task 5: 添加日志

**Files:**
- Create: `backend/internal/logger/logger.go`
- Modify: `backend/cmd/main.go`

- [ ] **Step 1: 创建日志模块**

```go
package logger

import (
    "log"
    "os"
    "time"
)

var (
    Info  *log.Logger
    Error *log.Logger
    Debug *log.Logger
)

func Init() {
    Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
    Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogRequest(method, path string, status, duration int) {
    Info.Printf("%s %s %d %dms", method, path, status, duration)
}
```

- [ ] **Step 2: 添加日志中间件**

```go
package middleware

import (
    "time"

    "life-financial-assistant-backend/internal/logger"
    "github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := int(time.Since(start).Milliseconds())
        logger.LogRequest(c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
    }
}
```

- [ ] **Step 3: 更新 main.go 启用日志**

```go
func main() {
    logger.Init()
    
    r := gin.New()
    r.Use(middleware.RecoveryMiddleware())
    r.Use(middleware.LoggerMiddleware())
    r.Use(middleware.CORSMiddleware())
    // ...
}
```

- [ ] **Step 4: 提交**

```bash
git add backend/internal/logger/ backend/internal/api/middleware/
git commit -m "feat: 添加请求日志中间件"
```

---

## 文件变更清单

| 操作 | 文件 |
|------|------|
| 新增 | `backend/internal/config/config.go` |
| 新增 | `backend/internal/error/errors.go` |
| 新增 | `backend/internal/dto/auth.go` |
| 新增 | `backend/internal/dto/space.go` |
| 新增 | `backend/internal/dto/bill.go` |
| 新增 | `backend/internal/logger/logger.go` |
| 新增 | `backend/internal/api/middleware/recovery.go` |
| 重命名 | `auth_middleware.go` → `auth.go` |
| 重命名 | `package middleware.go` → `cors.go` |
| 修改 | `backend/cmd/main.go` |
| 修改 | `backend/internal/api/handler/auth_handler.go` |