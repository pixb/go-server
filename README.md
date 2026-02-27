# go-server

练习建立标准的 `go-server` 服务端。

写一个用户登录的示例程序。

- 数据库使用 `sqlite`.

## proto定义

用户信息：

- id
- username
- nickname
- email
- phone
- role
- password_expires_at
- created_at
- updated_at

### 文件夹创建

协议文件创建在 `proto/api/v1` 目录下.

初始化 buf 工程。

```shell
cd proto
buf config init
```

生成 `buf.yaml` 文件。

编辑 `buf.yaml`.

加入 `googleapis` 引用.

```c
deps:
  - buf.build/googleapis/googleapis
```

执行

```shell
buf dep update
```

添加 `buf.gen.yaml` 生成配置。

### 创建用户协议

`common.proto`

### 生成配置

```shell
buf generate
```

### 定义注册用户的服务及消息对象

- 通用: common
  - 用户信息: User

- 用户服务UserService{}
  - 注册用户RegisterUser()
    - RegisterUserRequest
    - RegisterUserResponse
  - 获取用户信息: GetUserProfile()
    - GetUserProfileRequest
    - GetUserProfileResponse
  - 更新用户信息: UpdateUserProfile()
    - UpdateUserProfileRequest
    - UpdateUserProfileResponse
  - 修改密码: ChangePassword()
    - ChangePasswordRequest
    - ChangePasswordResponse

- 认证服务AuthService{}
  - 登录用户LoginUser()
    - LoginUserRequest
    - LoginUserResponse
  - 刷新Token: RefreshToken()
    - RefreshTokenRequest
    - RefreshTokenResponse
  - 验证Token: ValidateToken
    - ValidateTokenRequest
    - ValidateTokenResponse
  - 登出: Logout
    - LogoutRequest
    - LogoutResponse

重新生成。

```shell
buf generate
```

### 入口main

- `cmd/server/main.go`

读取配置信息，根据配置启动服务端。

定义流程

- 定义 `rootCmd`
- 定义 `init()`
- 定义 `main()` 方法

使用 `cobra` 库, 定义命令行参数，使用 `viper` 管理配置信息。

`rootCmd` 接收命令行参数。

- `Mode`: 模式。
- `Addr`: 服务器地址.
- `Port`: 端口
- `Data`: 数据目录
- `Driver`: 数据库驱动
- `DSN`: 数据库连接连接地址
- `Secret`: 密钥
- `Version`: 版本

### init()方法

设置属性的提示值和默认值.

### 创建Profile配置文件

- `internal/profile/profile.go`

### 定义run()方法，运行服务器

- 在 `rootCmd` 中解析配置文件.
- main.go 中定义 `run()` 方法.

`run()` 步骤:

- 1.检查配置是否正确
- 2.创建数据目录
- 3.创建数据驱动
- 4.创建存储实例并且迁移数据

## store 模块

store 就是存储模块，数据库访问相关。

### store.go

首先定义 `store/store.go` 文件。

- 定义 `Driver` 接口.
- 定义 `Store` 结构.
- 创建的方法 `New()`.
- 获取驱动方法 `GetDriver()`.
- 关闭的方法 `Close()`.
- 创建用户的方法 `CreateUser()`.
- 获取用户信息的方法 `GetUser()`.
- 获取用户列表 `ListUsers()`.
- 更新用户信息 `UpdateUser()`.
- 删除用户 `DeleteUser()`.
- 根据用户名获取用户信息 `GetUserByUsername()`.
- 根据用户Email获取用户信息 `GetUserByEmail()`.
- 创建刷新`Token`: `CreateRefreshToken()`.
- 更新刷新`Token`: `UpdateRefreshToken()`.
- 列出刷新`Token`: `ListRefreshTokens()`.
- 删除刷新`Token`: `DeleteRefreshToken()`.
- 获取刷新`Token`: `GetRefreshToken()`.
- 迁移方法 `Migrate()`.

### cache.go

缓存模块 `store/cache/cache.go`.

### 数据模型 model.go

创建数据库相关模型: `store/model.go`.

### 创建迁移文件

```c
migration
├── postgresql
│   ├── LATEST.sql
│   └── refresh_tokens.sql
└── sqlite
    ├── LATEST.sql
    └── refresh_tokens.sql

```

### 创建数据库操作实现

```c
db
├── postgresql
│   └── postgresql.go
└── sqlite
    └── sqlite.go

```

## 定义服务模块server

服务模块 `server/`.

### 创建服务入口主文件

- `server/server.go`
- 定义服务结构 `Server`.
- 定义创建服务的方法: `NewServer()`.
- 定义启动服务的方法: `Start()`.
- 定义关闭服务的方法: `Shutdown()`.

### 定义认证

创建 `server/auth` 路径。

#### token

创建 `server/auth/token.go`

- 创建自定义 `JWTClaims`.
- 创建生成Token的方法 `GenerateAccessToken()`.
- 创建验证Token的方法 `ValidateAccessToken()`.
- 创建生成刷新Token的方法 `GenerateRefreshToken()`.
- 创建生成密码的方法 `HashPassword()`.
- 创建校验密码的方法 `CheckPassword()`.

#### claims

创建 `server/auth/claims.go`

用于使用上下文来传递用户身份的信息。

- 创建 `UserClaims` 结构。
- 创建Context枚举key.
- 从上下文获取用户ID: `GetUserID()`.
- 从上下文获取用户声明信息: `GetUserClaims()`, `UserClaims`.
- 设置用户声明信息到上下文: `SetUserClaimsInContext()`.
- 设置用户信息到上下文：`SetUserInContext()`

#### extract

创建 `server/auth/extract.go`

- 定义从认证头中提取 `Bearer` 认证信息的方法.
- `ExtractBearerToken()`

#### authenticator

创建 `server/auth/authenticator.go`

- 定义认证器结构: `Authenticator`.
- 定义创建认证器函数: `NewAuthenticator()`.
- 定义认证函数: `Authenticate()`.

### 定义访问控制列表 acl_config.go

- `server/router/api/v1/acl_config.go`
- `PublicMethods`: 公有方法映射表
- 添加是否是共有方法的函数: `IsPublicMethod()`.

### 定义 `connectRPC` 的 Handler 类

- `server/router/api/v1/connect_handler.go`

### 定义 user_service.go 实现connectrpc需要的实现方法

- `server/router/api/v1/user_service.go`

### 定义 connect_interceptors.go 拦截器

- `server/router/api/v1/connect_interceptors.go`

### 定义日志拦截器 logging_interceptor.go

- `server/router/api/v1/logging_interceptor.go`

### 定义 recovery_interceptor.go

- `server/router/api/v1/recovery_interceptor.go`

### 定义 auth_interceptor.go

- `server/router/api/v1/auth_interceptor.go`

### 定义common 模块

- `server/common/`

#### 定义response.go

- `server/common/response.go`

### 定义v1服务文件

- `server/router/api/v1/v1.go`
- 定义v1服务结构: `APIV1Service`.
- 定义创建结构实例方法: `NewAPIV1Service()`
- 定义连接状态到状态码的转换函数: `connectCodeToState()`.
- 定义注册 `gateway` 方法: `RegisterGateway()`.

## 数据库迁移定义

- [迁移文件](store/migration)
- [seed示例数据库脚本](store/seed)
- [迁移代码](store/migrator.go)

