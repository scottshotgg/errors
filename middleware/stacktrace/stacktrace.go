package stacktrace

import (
	"context"

	"github.com/scottshotgg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryGRPC(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Run the handler
		// We will not check the error here until later
		var (
			resp, err  = handler(ctx, req)
			pkgErr, ok = err.(*errors.Err)
		)

		// TODO: figure this out
		if ok {
			metadata.AppendToOutgoingContext(ctx, "stacktrace", pkgErr.Stack().String())
		}

		// Pack stack trace into metadata

		return resp, err
	}

}
