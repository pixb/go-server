package middleware

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// UnifiedResponse represents the standard JSON response format for all API endpoints
type UnifiedResponse struct {
	State   int         `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewSuccessResponse creates a successful response
func NewSuccessResponse(data interface{}) UnifiedResponse {
	return UnifiedResponse{
		State:   0,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(state int, message string) UnifiedResponse {
	return UnifiedResponse{
		State:   state,
		Message: message,
		Data:    nil,
	}
}

// NewCORSHandler creates a CORS middleware with proper configuration
func NewCORSHandler() echo.MiddlewareFunc {
	return echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOriginFunc: func(_ string) (bool, error) {
			return true, nil
		},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodPatch},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	})
}

// UnifiedResponseMiddleware creates a middleware that统一处理所有响应格式
func UnifiedResponseMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 执行后续处理程序
			err := next(c)

			// 如果没有错误，直接返回
			if err == nil {
				// 检查响应是否已经被写入
				if c.Response().Committed {
					return nil
				}

				// 检查响应体是否为空
				if c.Response().Size == 0 {
					// 如果响应体为空，返回成功响应
					return c.JSON(http.StatusOK, NewSuccessResponse(nil))
				}

				return nil
			}

			// 处理 echo.HTTPError 类型的错误
			if httpErr, ok := err.(*echo.HTTPError); ok {
				statusCode := httpErr.Code
				message := httpErr.Message
				_ = c.JSON(statusCode, NewErrorResponse(statusCode, message.(string)))
				return err
			}

			// 处理不同类型的错误
			var statusCode int
			var state int
			var message string

			// 处理 Connect 错误
			if connectErr, ok := err.(*connect.Error); ok {
				message = connectErr.Message()
				state = connectCodeToState(connectErr.Code().String())

				// 映射 Connect 错误码到 HTTP 状态码
				switch connectErr.Code() {
				case connect.CodeCanceled:
					statusCode = http.StatusRequestTimeout
				case connect.CodeUnknown:
					statusCode = http.StatusInternalServerError
				case connect.CodeInvalidArgument:
					statusCode = http.StatusBadRequest
				case connect.CodeDeadlineExceeded:
					statusCode = http.StatusRequestTimeout
				case connect.CodeNotFound:
					statusCode = http.StatusNotFound
				case connect.CodeAlreadyExists:
					statusCode = http.StatusConflict
				case connect.CodePermissionDenied:
					statusCode = http.StatusForbidden
				case connect.CodeUnauthenticated:
					statusCode = http.StatusUnauthorized
				case connect.CodeResourceExhausted:
					statusCode = http.StatusTooManyRequests
				case connect.CodeFailedPrecondition:
					statusCode = http.StatusBadRequest
				case connect.CodeAborted:
					statusCode = http.StatusConflict
				case connect.CodeOutOfRange:
					statusCode = http.StatusBadRequest
				case connect.CodeUnimplemented:
					statusCode = http.StatusNotImplemented
				case connect.CodeInternal:
					statusCode = http.StatusInternalServerError
				case connect.CodeUnavailable:
					statusCode = http.StatusServiceUnavailable
				case connect.CodeDataLoss:
					statusCode = http.StatusInternalServerError
				default:
					statusCode = http.StatusInternalServerError
				}
			} else {
				// 处理其他类型的错误
				message = err.Error()
				state = 2 // 通用错误状态码
				statusCode = http.StatusInternalServerError
			}

			// 返回统一格式的错误响应
			return c.JSON(statusCode, NewErrorResponse(state, message))
		}
	}
}

// connectCodeToState 将 Connect 错误码转换为统一状态码
func connectCodeToState(code string) int {
	switch code {
	case "CANCELLED":
		return 1
	case "UNKNOWN":
		return 2
	case "INVALID_ARGUMENT":
		return 3
	case "DEADLINE_EXCEEDED":
		return 4
	case "NOT_FOUND":
		return 5
	case "ALREADY_EXISTS":
		return 6
	case "PERMISSION_DENIED":
		return 7
	case "UNAUTHENTICATED":
		return 16
	case "RESOURCE_EXHAUSTED":
		return 8
	case "FAILED_PRECONDITION":
		return 9
	case "ABORTED":
		return 10
	case "OUT_OF_RANGE":
		return 11
	case "UNIMPLEMENTED":
		return 12
	case "INTERNAL":
		return 13
	case "UNAVAILABLE":
		return 14
	case "DATA_LOSS":
		return 15
	default:
		return 2 // UNKNOWN
	}
}
