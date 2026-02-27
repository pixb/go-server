package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewGatewayErrorHandler creates a custom error handler for gRPC-Gateway
func NewGatewayErrorHandler() func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	return func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		// Convert gRPC error to unified response format
		var statusCode int
		var state int
		var message string

		// Default values
		statusCode = http.StatusInternalServerError
		state = 2
		message = "internal server error"

		// Check if it's a gRPC status error
		if grpcErr, ok := status.FromError(err); ok {
			message = grpcErr.Message()
			state = int(grpcErr.Code())

			// Map gRPC codes to HTTP status codes
			switch grpcErr.Code() {
			case codes.Unauthenticated:
				statusCode = http.StatusUnauthorized
			case codes.InvalidArgument:
				statusCode = http.StatusBadRequest
			case codes.NotFound:
				statusCode = http.StatusNotFound
			case codes.AlreadyExists:
				statusCode = http.StatusConflict
			default:
				statusCode = http.StatusInternalServerError
			}
		} else if connectErr, ok := err.(*connect.Error); ok {
			// Handle Connect errors
			message = connectErr.Message()
			state = connectCodeToState(connectErr.Code().String())

			// Map Connect codes to HTTP status codes
			switch connectErr.Code() {
			case connect.CodeUnauthenticated:
				statusCode = http.StatusUnauthorized
			case connect.CodeInvalidArgument:
				statusCode = http.StatusBadRequest
			case connect.CodeNotFound:
				statusCode = http.StatusNotFound
			case connect.CodeAlreadyExists:
				statusCode = http.StatusConflict
			default:
				statusCode = http.StatusInternalServerError
			}
		} else {
			// Handle other errors
			message = err.Error()
			state = 2
			statusCode = http.StatusInternalServerError
		}

		// Write unified error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		// Create unified response
		unifiedResp := NewErrorResponse(state, message)

		// Encode and write response
		json.NewEncoder(w).Encode(unifiedResp)
	}
}

// NewGatewayResponseWrapper creates a response wrapper for gRPC-Gateway
func NewGatewayResponseWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer that intercepts the response
		buffer := &bytes.Buffer{}

		// Create a custom response writer that captures the response
		customWriter := &responseWriter{
			originalWriter: w,
			buffer:         buffer,
			statusCode:     http.StatusOK,
		}

		// Call the next handler with the custom writer
		next.ServeHTTP(customWriter, r)

		// If response has been written, wrap it in unified format
		if buffer.Len() > 0 {
			// Check if response is JSON
			if contentType := customWriter.Header().Get("Content-Type"); contentType == "application/json" {
				// Read the original response
				originalResponse := buffer.Bytes()

				// Check if it's already a unified response
				var unifiedResp UnifiedResponse
				if err := json.Unmarshal(originalResponse, &unifiedResp); err == nil && unifiedResp.State != 0 {
					// It's already an error response, just write it
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(customWriter.statusCode)
					w.Write(originalResponse)
					return
				}

				// It's a success response, wrap it in unified format
				// Reset the response writer
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(customWriter.statusCode)

				// Create a unified response
				var data interface{}
				if err := json.Unmarshal(originalResponse, &data); err == nil {
					unifiedResp := NewSuccessResponse(data)
					json.NewEncoder(w).Encode(unifiedResp)
				} else {
					// If unmarshaling fails, just write the original response
					w.Write(originalResponse)
				}
			} else {
				// If response is not JSON, just write the original response
				w.Write(buffer.Bytes())
			}
		}
	})
}

// responseWriter is a custom http.ResponseWriter that captures the response
type responseWriter struct {
	originalWriter http.ResponseWriter
	buffer         *bytes.Buffer
	statusCode     int
}

// WriteHeader captures the WriteHeader method
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	// Don't write header to the original response yet
}

// Write captures the Write method
func (w *responseWriter) Write(b []byte) (int, error) {
	return w.buffer.Write(b)
}

// Header returns the header map
func (w *responseWriter) Header() http.Header {
	return w.originalWriter.Header()
}
