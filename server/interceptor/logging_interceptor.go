package interceptor

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
)

type LoggingInterceptor struct {
	logStacktraces bool
}

func NewLoggingInterceptor(logStacktraces bool) *LoggingInterceptor {
	return &LoggingInterceptor{
		logStacktraces: logStacktraces,
	}
}

func (in *LoggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()

		resp, err := next(ctx, req)

		duration := time.Since(start)
		status := "OK"
		if err != nil {
			status = "ERROR"
		}

		fmt.Printf("[connect] %s %s %v\n", status, req.Spec().Procedure, duration)

		if err != nil && in.logStacktraces {
			fmt.Printf("[connect] Error: %v\n", err)
		}

		return resp, err
	}
}

func (in *LoggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (in *LoggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
