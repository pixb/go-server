package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/grpc/metadata"
)

type MetadataInterceptor struct{}

func NewMetadataInterceptor() *MetadataInterceptor {
	return &MetadataInterceptor{}
}

func (*MetadataInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		header := req.Header()
		md := metadata.MD{}

		if ua := header.Get("User-Agent"); ua != "" {
			md.Set("user-agent", ua)
		}
		if xff := header.Get("X-Forwarded-For"); xff != "" {
			md.Set("x-forwarded-for", xff)
		}
		if xri := header.Get("X-Real-Ip"); xri != "" {
			md.Set("x-real-ip", xri)
		}
		if cookie := header.Get("Cookie"); cookie != "" {
			md.Set("cookie", cookie)
		}

		ctx = metadata.NewIncomingContext(ctx, md)

		resp, err := next(ctx, req)

		// Only set headers if there's no error and resp is not nil
		if err == nil && resp != nil {
			header := resp.Header()
			if header != nil {
				header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
				header.Set("Pragma", "no-cache")
				header.Set("Expires", "0")
			}
		}

		return resp, err
	}
}

func (*MetadataInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (*MetadataInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
