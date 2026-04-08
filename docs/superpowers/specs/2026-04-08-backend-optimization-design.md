# 后端优化设计方案

**日期**: 2026-04-08  
**主题**: Go 后端架构优化  
**当前状态**: 已审查，待实现

---

## 1. 代码结构优化

### 1.1 当前问题

- `main.go` 混合了初始化、路由注册、服务启动 (~80行)
- service 层未进行接口抽象，直接依赖具体 repository
- 缺少统一的错误类型定义

### 1.2 优化方案

| 文件 | 优化内容 |
|------|----------|
| `cmd/main.go` | 仅保留入口逻辑，提取初始化到 `internal/bootstrap/` |
| `internal/config/` | 新增配置模块，集中管理配置 |
| `internal/service/interfaces.go` | 新增服务接口定义 |
| `internal/error/` | 新增自定义错误类型 |

**目录结构**:

```
backend/internal/
├── bootstrap/          # 新增：初始化逻辑
│   └── init.go
├── config/              # 新增：配置管理
│   └── config.go
├── error/              # 新增：自定义错误
│   └── errors.go
├── service/
│   └── interfaces.go   # 新增：服务接口
├── ...
```

---

## 2. 错误处理统一

### 2.1 当前问题

- 各 handler 直接返回 `gin.H{"error": "..."}` 
- 错误格式不统一
- 缺少错误码

### 2.2 优化方案

```go
// internal/error/errors.go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func NewBadRequestError(msg string) *AppError {
    return &AppError{Code: 400, Message: msg}
}

func NewUnauthorizedError(msg string) *AppError {
    return &AppError{Code: 401, Message: msg}
}

func NewInternalError(msg string) *AppError {
    return &AppError{Code: 500, Message: msg}
}
```

```go
// 统一响应助手
func RespondError(c *gin.Context, err *AppError) {
    c.JSON(err.Code, gin.H{"error": err.Message, "code": err.Code})
}
```

**优化后 handler 示例**:

```go
func (h *AuthHandler) Login(c *gin.Context) {
    if err := c.ShouldBindJSON(&req); err != nil {
        RespondError(c, NewBadRequestError("invalid request body"))
        return
    }
    // ...
}
```

---

## 3. 中间件整理

### 3.1 当前问题

- `auth_middleware.go` + `middleware.go` 两个文件，命名不规范（空格）
- CORS 使用 `*`，不安全

### 3.2 优化方案

| 文件 | 优化内容 |
|------|----------|
| `middleware/auth.go` | 重命名，原地优化 |
| `middleware/cors.go` | 重命名 |
| `middleware/recovery.go` | 新增：panic 恢复中间件 |
| `middleware/logger.go` | 新增：请求日志中间件 |

**CORS 安全化**:

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        // ... 限定允许的域
    }
}
```

---

## 4. 类型/模型完善

### 4.1 当前问题

- `user_space.go` 混合 User、Space、SpaceUserLink
- 缺少 DTO/Request/Response 定义
- model 直接暴露 JSON tag

### 4.2 优化方案

```
backend/internal/
├── dto/                    # 新增：数据传输对象
│   ├── auth.go
│   ├── space.go
│   └── bill.go
└── model/
    ├── user.go            # 拆分：仅保留 User
    ├── space.go          # 拆分：仅保留 Space
    ├── bill.go          # 拆分：仅保留 Bill
    └── role.go
```

**DTO 示例**:

```go
// internal/dto/auth.go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
}

type LoginResponse struct {
    Token string `json:"token"`
}
```

---

## 5. 添加日志/监控

### 5.1 依赖

```bash
go get -u github.com/gin-contrib/zap
go get -u github.com/go-playground/validator/v10
```

### 5.2 优化方案

| 文件 | 内容 |
|------|------|
| `middleware/logger.go` | 使用 zap 记录请求日志 |
| `internal/logger/logger.go` | 统一日志配置 |

**日志格式**:

```
2026/04/08 10:30:45 [INFO] POST /api/bills 200 15ms
2026/04/08 10:30:46 [ERROR] POST /api/bills 400 5ms - invalid request body
```

---

## 6. 环境变量配置

### 6.1 当前问题

- API Key 硬编码：`AIAgent{APIKey: ""}`
- 数据库路径硬编码
- 缺少环境区分

### 6.2 优化方案

```go
// internal/config/config.go
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

func Load() (*Config, error) {
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
    }, nil
}
```

---

## 文件变更清单

| 操作 | 文件路径 |
|------|----------|
| 新增 | `internal/bootstrap/init.go` |
| 新增 | `internal/config/config.go` |
| 新增 | `internal/error/errors.go` |
| 新增 | `internal/dto/auth.go` |
| 新增 | `internal/dto/space.go` |
| 新增 | `internal/dto/bill.go` |
| 新增 | `internal/logger/logger.go` |
| 新增 | `middleware/logger.go` |
| 新增 | `middleware/recovery.go` |
| 重命名 | `middleware/auth_middleware.go` → `middleware/auth.go` |
| 重命名 | `middleware/package middleware.go` → `middleware/cors.go` |
| 修改 | `cmd/main.go` |
| 修改 | `internal/model/user_space.go` → 拆分 |

---

## 实施顺序

1. **环境变量配置** — 为其他改动打基础
2. **错误处理统一** — 定义错误类型
3. **中间件整理** — 规范化文件名
4. **类型/模型完善** — 拆分 model + 添加 DTO
5. **代码结构优化** — 提取初始化逻辑
6. **添加日志** — 结构化日志

---

## 待确认

- [ ] 日志级别配置（INFO/DEBUG）是否需要？
- [ ] 生产环境是否需要区分？
