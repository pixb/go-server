package auth

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pixb/go-server/server/common"
)

// NewGatewayAuthMiddleware creates a gRPC-Gateway authentication middleware
func NewGatewayAuthMiddleware(authenticator *Authenticator) func(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			ctx := r.Context()

			// Get the RPC method name from context
			rpcMethod, ok := runtime.RPCMethod(ctx)

			// Extract credentials from HTTP headers
			authHeader := r.Header.Get("Authorization")

			// Execute authentication
			result := authenticator.Authenticate(ctx, authHeader)

			// Enforce authentication for non-public methods
			if result == nil && ok && !common.IsPublicMethod(rpcMethod) {
				http.Error(w, `{"state": 401, "message": "authentication required", "data": null}`, http.StatusUnauthorized)
				return
			}

			// Set context based on auth result
			if result != nil {
				if result.Claims != nil {
					ctx = SetUserClaimsInContext(ctx, result.Claims)
					ctx = SetUserIDInContext(ctx, result.Claims.UserID)
				} else if result.User != nil {
					ctx = SetUserInContext(ctx, result.User, result.AccessToken)
				}
				r = r.WithContext(ctx)
			}

			next(w, r, pathParams)
		}
	}
}

// SetUserIDInContext sets the user ID in the context
func SetUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
