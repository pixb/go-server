package v1

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	v1connect "github.com/pixb/go-server/proto/gen/api/v1/apiv1connect"
)

type ConnectServiceHandler struct {
	*APIV1Service
}

func NewConnectServiceHandler(svc *APIV1Service) *ConnectServiceHandler {
	return &ConnectServiceHandler{APIV1Service: svc}
}

func (s *ConnectServiceHandler) RegisterConnectHandlers(mux *http.ServeMux, opts ...connect.HandlerOption) {
	// Register UserService handler
	userPath, userHandler := v1connect.NewUserServiceHandler(s, opts...)
	mux.Handle(userPath, userHandler)

	// Register AuthService handler
	authPath, authHandler := v1connect.NewAuthServiceHandler(s, opts...)
	mux.Handle(authPath, authHandler)
}

func (s *ConnectServiceHandler) RegisterUser(ctx context.Context, req *connect.Request[v1pb.RegisterUserRequest]) (*connect.Response[v1pb.RegisterUserResponse], error) {
	resp, err := s.APIV1Service.RegisterUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) Login(ctx context.Context, req *connect.Request[v1pb.LoginRequest]) (*connect.Response[v1pb.LoginResponse], error) {
	resp, err := s.APIV1Service.Login(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) RefreshToken(ctx context.Context, req *connect.Request[v1pb.RefreshTokenRequest]) (*connect.Response[v1pb.RefreshTokenResponse], error) {
	resp, err := s.APIV1Service.RefreshToken(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) ValidateToken(ctx context.Context, req *connect.Request[v1pb.ValidateTokenRequest]) (*connect.Response[v1pb.ValidateTokenResponse], error) {
	resp, err := s.APIV1Service.ValidateToken(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) Logout(ctx context.Context, req *connect.Request[v1pb.LogoutRequest]) (*connect.Response[v1pb.LogoutResponse], error) {
	resp, err := s.APIV1Service.Logout(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) ChangePassword(ctx context.Context, req *connect.Request[v1pb.ChangePasswordRequest]) (*connect.Response[v1pb.ChangePasswordResponse], error) {
	resp, err := s.APIV1Service.ChangePassword(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) GetUserProfile(ctx context.Context, req *connect.Request[v1pb.GetUserProfileRequest]) (*connect.Response[v1pb.GetUserProfileResponse], error) {
	resp, err := s.APIV1Service.GetUserProfile(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *ConnectServiceHandler) UpdateUserProfile(ctx context.Context, req *connect.Request[v1pb.UpdateUserProfileRequest]) (*connect.Response[v1pb.UpdateUserProfileResponse], error) {
	resp, err := s.APIV1Service.UpdateUserProfile(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
