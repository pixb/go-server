package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pixb/go-server/internal/profile"
	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/server/auth"
	"github.com/pixb/go-server/server/common"
	"github.com/pixb/go-server/server/middleware"
	v1 "github.com/pixb/go-server/server/router/api/v1"
	"github.com/pixb/go-server/store"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// 1.定义服务结构
type Server struct {
	Profile *profile.Profile
	Store   *store.Store
	Secret  string

	echoServer         *echo.Echo
	grpcServer         *grpc.Server
	apiV1Service       *v1.APIV1Service
	healthCheckService *common.HealthCheckService
	wg                 sync.WaitGroup
}

// 2.创建服务实例指针的方法
func NewServer(ctx context.Context, prof *profile.Profile, store *store.Store) (*Server, error) {
	s := &Server{
		Profile: prof,
		Store:   store,
	}

	// 2.1. 创建Echo服务实例
	echoServer := echo.New()
	echoServer.Debug = prof.Demo
	echoServer.HideBanner = true
	echoServer.Use(echomiddleware.Recover())
	echoServer.Use(echomiddleware.Logger())
	echoServer.Use(middleware.NewCORSHandler())

	// Only enable rate limiter in production mode
	// if !prof.IsDev() {
	// 	echoServer.Use(middleware.NewRateLimiter(middleware.DefaultRateLimiterConfig()))
	// }

	// Only enable CSRF protection in production mode
	// if !prof.IsDev() {
	// 	echoServer.Use(middleware.CSRFTokenMiddleware(middleware.DefaultCSRFConfig()))
	// 	echoServer.Use(middleware.NewCSRF(middleware.DefaultCSRFConfig()))
	// }

	s.echoServer = echoServer

	// Initialize health check service
	dbChecker := common.NewDatabaseChecker(store)
	serviceChecker := common.NewServiceChecker()
	s.healthCheckService = common.NewHealthCheckService(dbChecker, serviceChecker)

	// Register health check endpoints
	echoServer.GET("/healthz", common.HealthCheckHandler(s.healthCheckService))
	echoServer.GET("/readyz", func(c echo.Context) error {
		ctx := c.Request().Context()
		if s.healthCheckService.IsHealthy(ctx) {
			return c.String(http.StatusOK, "Service ready.")
		}
		return c.String(http.StatusServiceUnavailable, "Service not ready.")
	})

	if prof.Secret == "" {
		prof.Secret = "your-secret-key"
	}
	s.Secret = prof.Secret

	s.apiV1Service = v1.NewAPIV1Service(s.Secret, prof, store)

	authInterceptor := auth.NewInterceptor(store, s.Secret)
	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.GRPCUnaryInterceptor()))
	v1pb.RegisterUserServiceServer(s.grpcServer, s.apiV1Service)
	v1pb.RegisterAuthServiceServer(s.grpcServer, s.apiV1Service)

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.Profile.Addr, s.Profile.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	m := cmux.New(listener)
	httpListener := m.Match(cmux.HTTP1Fast())
	grpcListener := m.Match(cmux.HTTP2())

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.echoServer.Logger.Info("gRPC server starting on ", address)
		s.grpcServer.Serve(grpcListener)
	}()

	if err := s.apiV1Service.RegisterGateway(ctx, s.echoServer); err != nil {
		return err
	}

	s.echoServer.Listener = httpListener
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.echoServer.Logger.Info("Echo server starting on ", address)
		s.echoServer.Start(address)
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		m.Serve()
	}()

	s.echoServer.Logger.Info("Server started successfully (HTTP/1.1 + HTTP/2)")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.echoServer.Logger.Info("Shutting down server...")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	s.echoServer.Shutdown(ctx)

	s.Store.Close()
	s.wg.Wait()
	return nil
}
