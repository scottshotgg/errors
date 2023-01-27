package logging

import (
	"context"
	"errors"
	"time"

	pkgerrors "github.com/scottshotgg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Option int

const (
	_ Option = iota

	DisableStackTrace
	DisableCause
	DisableError
)

type GRPCLogger struct {
	logger *zap.Logger

	disableCause      bool
	disableStackTrace bool
	disableError      bool
}

func New(logger *zap.Logger, opts ...Option) (*GRPCLogger, error) {
	if logger == nil {
		return nil, errors.New("nil logger")
	}

	var g = GRPCLogger{
		logger: logger,
	}

	for _, opt := range opts {
		switch opt {
		case DisableCause:
			g.disableCause = true

		case DisableStackTrace:
			g.disableStackTrace = true

		case DisableError:
			g.disableError = true

		}
	}

	return &g, nil
}

func (g *GRPCLogger) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		// Catch the start time
		start = time.Now()

		// Run the handler
		// We will not check the error here until later
		resp, err = handler(ctx, req)
	)

	g.logger.Debug("finished unary call",
		g.fields(start, err)...,
	)

	// Pack stack trace into metadata

	return resp, err
}

func (g *GRPCLogger) fields(start time.Time, err error, opts ...Option) []zap.Field {
	var (
		// Convert the error to a status; this will always work since GRPC always returns this error
		st = status.Convert(err)

		code = st.Code()

		// Make the base fields
		fields = []zap.Field{
			zap.String("start_time", start.Format(time.RFC3339)),
			zap.String("duration", time.Since(start).String()),
			zap.Int("status.code", int(code)),
			zap.String("status.code.desc", code.String()),
		}
	)

	var msg = st.Message()
	if msg != "" {
		fields = append(fields,
			zap.String("status.message", st.Message()),
		)
	}

	// If there are details then append those
	if len(st.Details()) > 0 {
		fields = append(fields,
			zap.Any("status.details", st.Details()),
		)
	}

	if g.disableError {
		return fields
	}

	// If the error is non-nil then apply the error fields
	if err != nil {
		var pkgErr, ok = err.(*pkgerrors.Err)
		if !ok {
			return fields
		}

		if pkgErr == nil {
			return fields
		}

		fields = append(fields,
			zap.String("error", pkgErr.Error()),
		)

		if !g.disableCause {
			var cause = pkgErr.Cause()
			if cause != nil {
				fields = append(fields,
					zap.String("cause.name", cause.Name()),
					zap.String("cause.error", cause.Error()),
				)
			}
		}

		if !g.disableStackTrace {
			var stack = pkgErr.Stack()

			// Append the stacktrace
			if stack != nil {
				fields = append(fields,
					zap.Any("stacktrace", stack.Entries()),
					// zap.String("stacktrace", pkgErr.StackString()),
				)
			}
		}
	}

	return fields
}
