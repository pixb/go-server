package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pixb/echo-demo/go-server/store"
)

// 健康检查服务定义
// ① 定义健康检查接口
// ② 定义健康检查服务结构
// ③ 定义创建健康检查服务函数
// ④ 定义健康检查服务检查方法
// ⑤ 定义健康检查服务是否健康方法

// HealthChecker is an interface for health checkers
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

// HealthCheckService is a service that checks the health of the application
type HealthCheckService struct {
	checkers []HealthChecker
}

// NewHealthCheckService creates a new health check service
func NewHealthCheckService(checkers ...HealthChecker) *HealthCheckService {
	return &HealthCheckService{
		checkers: checkers,
	}
}

// Check checks the health of all registered checkers
func (s *HealthCheckService) Check(ctx context.Context) map[string]error {
	results := make(map[string]error)
	for _, checker := range s.checkers {
		results[checker.Name()] = checker.Check(ctx)
	}
	return results
}

// IsHealthy checks if all registered checkers are healthy
func (s *HealthCheckService) IsHealthy(ctx context.Context) bool {
	results := s.Check(ctx)
	for _, err := range results {
		if err != nil {
			return false
		}
	}
	return true
}

// 数据库检查器定义
// ⑥ 定义数据库检查器结构
// ⑦ 定义创建数据库检查器函数
// ⑧ 定义数据库检查器检查方法
// ⑨ 定义数据库检查器名称方法

// DatabaseChecker is a health checker for database connections
type DatabaseChecker struct {
	store *store.Store
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(store *store.Store) *DatabaseChecker {
	return &DatabaseChecker{
		store: store,
	}
}

// Check checks the health of the database connection
func (c *DatabaseChecker) Check(ctx context.Context) error {
	if c.store == nil {
		return errors.New("database store is nil")
	}

	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if the database is reachable
	if err := c.store.Ping(timeoutCtx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Name returns the name of the checker
func (c *DatabaseChecker) Name() string {
	return "database"
}

// ServiceChecker is a health checker for the service itself
type ServiceChecker struct {
	startTime time.Time
}

// NewServiceChecker creates a new service health checker
func NewServiceChecker() *ServiceChecker {
	return &ServiceChecker{
		startTime: time.Now(),
	}
}

// Check checks the health of the service
func (c *ServiceChecker) Check(ctx context.Context) error {
	// Service is always healthy if it's running
	return nil
}

// Name returns the name of the checker
func (c *ServiceChecker) Name() string {
	return "service"
}

// Uptime returns the uptime of the service
func (c *ServiceChecker) Uptime() time.Duration {
	return time.Since(c.startTime)
}

// HealthCheckHandler is a handler for health check requests
func HealthCheckHandler(service *HealthCheckService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		results := service.Check(ctx)
		isHealthy := service.IsHealthy(ctx)

		// Build response
		response := map[string]interface{}{
			"status":  "healthy",
			"checks":  make(map[string]interface{}),
			"version": "1.0.0",
		}

		if !isHealthy {
			response["status"] = "unhealthy"
		}

		// Add checker results
		checks := response["checks"].(map[string]interface{})
		for name, err := range results {
			checkStatus := "healthy"
			checkMessage := ""
			if err != nil {
				checkStatus = "unhealthy"
				checkMessage = err.Error()
			}
			checks[name] = map[string]string{
				"status":  checkStatus,
				"message": checkMessage,
			}
		}

		// Set status code
		statusCode := http.StatusOK
		if !isHealthy {
			statusCode = http.StatusServiceUnavailable
		}

		return c.JSON(statusCode, response)
	}
}
