package auth

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/pixb/echo-demo/go-server/server/common"
	"github.com/pixb/echo-demo/go-server/store"
)

// Interceptor is a common authentication interceptor that can be used for both gRPC and Connect-Go
type Interceptor struct {
	authenticator *Authenticator
}

// NewInterceptor creates a new authentication interceptor
func NewInterceptor(store *store.Store, secret string) *Interceptor {
	return &Interceptor{
		authenticator: NewAuthenticator(store, secret),
	}
}

// GRPCUnaryInterceptor returns a gRPC unary interceptor
func (i *Interceptor) GRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Extract authorization header
		authHeader := ""
		if vals := md.Get("authorization"); len(vals) > 0 {
			authHeader = vals[0]
		}

		// Execute authentication
		result := i.authenticator.Authenticate(ctx, authHeader)

		// Enforce authentication for non-public methods
		if result == nil && !common.IsPublicMethod(info.FullMethod) {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		// Set context based on auth result
		if result != nil {
			if result.Claims != nil {
				ctx = SetUserClaimsInContext(ctx, result.Claims)
				ctx = context.WithValue(ctx, UserIDContextKey, result.Claims.UserID)
			} else if result.User != nil {
				ctx = SetUserInContext(ctx, result.User, result.AccessToken)
			}
		}

		// Call the handler
		return handler(ctx, req)
	}
}

// ConnectUnaryInterceptor returns a Connect-Go unary interceptor
func (i *Interceptor) ConnectUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Extract authorization header
			authHeader := req.Header().Get("Authorization")

			// Execute authentication
			result := i.authenticator.Authenticate(ctx, authHeader)

			// Enforce authentication for non-public methods
			if result == nil && !common.IsPublicMethod(req.Spec().Procedure) {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
			}

			// Set context based on auth result
			if result != nil {
				if result.Claims != nil {
					ctx = SetUserClaimsInContext(ctx, result.Claims)
					ctx = context.WithValue(ctx, UserIDContextKey, result.Claims.UserID)
				} else if result.User != nil {
					ctx = SetUserInContext(ctx, result.User, result.AccessToken)
				}
			}

			// Call the handler
			return next(ctx, req)
		}
	}
}
