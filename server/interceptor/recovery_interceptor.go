package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	"connectrpc.com/connect"
)

type RecoveryInterceptor struct {
	logStacktraces bool
}

func NewRecoveryInterceptor(logStacktraces bool) *RecoveryInterceptor {
	return &RecoveryInterceptor{
		logStacktraces: logStacktraces,
	}
}

func (in *RecoveryInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		defer func() {
			if r := recover(); r != nil {
				if in.logStacktraces {
					fmt.Printf("[connect] Panic recovered: %v\n%s\n", r, debug.Stack())
				}
			}
		}()

		return next(ctx, req)
	}
}

func (in *RecoveryInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (in *RecoveryInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
