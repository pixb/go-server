package v1

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/pixb/go-server/internal/profile"
	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/server/auth"
	"github.com/pixb/go-server/server/interceptor"
	"github.com/pixb/go-server/server/middleware"
	"github.com/pixb/go-server/server/service"
	"github.com/pixb/go-server/store"
)

type APIV1Service struct {
	v1pb.UnimplementedUserServiceServer
	v1pb.UnimplementedAuthServiceServer

	Secret      string
	Profile     *profile.Profile
	Store       *store.Store
	UserService *service.UserService
	AuthService *service.AuthService
}

func NewAPIV1Service(secret string, profile *profile.Profile, store *store.Store) *APIV1Service {
	userService := service.NewUserService(secret, store)
	authService := service.NewAuthService(secret, store)
	return &APIV1Service{
		Secret:      secret,
		Profile:     profile,
		Store:       store,
		UserService: userService,
		AuthService: authService,
	}
}

// createConnectInterceptors creates Connect-Go interceptors
func (s *APIV1Service) createConnectInterceptors() connect.HandlerOption {
	logStacktraces := s.Profile.IsDev()
	authInterceptor := auth.NewInterceptor(s.Store, s.Secret)
	return connect.WithInterceptors(
		interceptor.NewMetadataInterceptor(),
		interceptor.NewLoggingInterceptor(logStacktraces),
		interceptor.NewRecoveryInterceptor(logStacktraces),
		authInterceptor.ConnectUnaryInterceptor(),
	)
}

func (s *APIV1Service) RegisterGateway(ctx context.Context, echoServer *echo.Echo) error {
	// =====================================================
	// STEP 1: Create Authenticator
	// =====================================================
	authenticator := auth.NewAuthenticator(s.Store, s.Secret)

	// =====================================================
	// STEP 2: Create gRPC-Gateway mux
	// =====================================================
	gwMux := runtime.NewServeMux(
		runtime.WithMiddlewares(auth.NewGatewayAuthMiddleware(authenticator)),
		runtime.WithErrorHandler(middleware.NewGatewayErrorHandler()),
	)

	// =====================================================
	// STEP 3: Register service handlers
	// =====================================================
	if err := v1pb.RegisterUserServiceHandlerServer(ctx, gwMux, s); err != nil {
		return err
	}
	if err := v1pb.RegisterAuthServiceHandlerServer(ctx, gwMux, s); err != nil {
		return err
	}

	// =====================================================
	// STEP 4: Create Connect service handler
	// =====================================================
	connectHandler := NewConnectServiceHandler(s)

	// =====================================================
	// STEP 5: Create common middlewares
	// =====================================================
	commonMiddlewares := []echo.MiddlewareFunc{
		middleware.NewCORSHandler(),
		middleware.UnifiedResponseMiddleware(),
	}

	// =====================================================
	// STEP 6: Register Connect handlers
	// =====================================================
	connectMux := http.NewServeMux()
	connectHandler.RegisterConnectHandlers(connectMux, s.createConnectInterceptors())

	// =====================================================
	// STEP 7: Register routes
	// =====================================================
	// Create gateway route group
	gwGroup := echoServer.Group("", commonMiddlewares...)
	gwGroup.Any("/api/v1/*", echo.WrapHandler(middleware.NewGatewayResponseWrapper(gwMux)))

	// Create connect route group
	connectGroup := echoServer.Group("", commonMiddlewares...)
	connectGroup.Any("/goserver.api.v1.*", echo.WrapHandler(connectMux))

	return nil
}

// UserService methods
func (s *APIV1Service) RegisterUser(ctx context.Context, req *v1pb.RegisterUserRequest) (*v1pb.RegisterUserResponse, error) {
	return s.UserService.RegisterUser(ctx, req)
}

func (s *APIV1Service) Login(ctx context.Context, req *v1pb.LoginRequest) (*v1pb.LoginResponse, error) {
	return s.AuthService.Login(ctx, req)
}

func (s *APIV1Service) RefreshToken(ctx context.Context, req *v1pb.RefreshTokenRequest) (*v1pb.RefreshTokenResponse, error) {
	return s.AuthService.RefreshToken(ctx, req)
}

func (s *APIV1Service) ValidateToken(ctx context.Context, req *v1pb.ValidateTokenRequest) (*v1pb.ValidateTokenResponse, error) {
	return s.AuthService.ValidateToken(ctx, req)
}

func (s *APIV1Service) Logout(ctx context.Context, req *v1pb.LogoutRequest) (*v1pb.LogoutResponse, error) {
	return s.AuthService.Logout(ctx, req)
}

func (s *APIV1Service) GetUserProfile(ctx context.Context, req *v1pb.GetUserProfileRequest) (*v1pb.GetUserProfileResponse, error) {
	return s.UserService.GetUserProfile(ctx, req)
}

func (s *APIV1Service) UpdateUserProfile(ctx context.Context, req *v1pb.UpdateUserProfileRequest) (*v1pb.UpdateUserProfileResponse, error) {
	return s.UserService.UpdateUserProfile(ctx, req)
}

func (s *APIV1Service) ChangePassword(ctx context.Context, req *v1pb.ChangePasswordRequest) (*v1pb.ChangePasswordResponse, error) {
	return s.UserService.ChangePassword(ctx, req)
}
