package logging

import (
	"component-master/config"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var counter uint64

type LogEntry struct {
	Level      string    `json:"level"`
	Timestamp  time.Time `json:"timestamp"`
	Message    string    `json:"message"`
	Error      string    `json:"error,omitempty"`
	Caller     string    `json:"caller,omitempty"`
	Additional any       `json:"additional,omitempty"`
}

type logHandler struct {
	slog.Handler
	additionalAttrs []slog.Attr
}

func InitLogger(cfg config.LogConfig) *slog.Logger {
	// Set default program level
	programLevel := new(slog.LevelVar)
	programLevel.Set(cfg.LogLevel)

	// Create handler options
	handlerOptions := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     programLevel,
	}

	// Create base handler based on output format
	var baseHandler slog.Handler
	if cfg.JSONOutput {
		baseHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	} else {
		baseHandler = slog.NewTextHandler(os.Stdout, handlerOptions)
	}

	// Create custom handler with additional attributes
	handler := &logHandler{
		Handler: baseHandler,
		additionalAttrs: []slog.Attr{
			slog.String("environment", cfg.Environment),
			slog.Time("boot_time", time.Now()),
		},
	}

	// Create and return the logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// Handle implements slog.Handler
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add additional attributes to the record
	for _, attr := range h.additionalAttrs {
		r.AddAttrs(attr)
	}

	return h.Handler.Handle(ctx, r)
}

// WithAttrs implements slog.Handler
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{
		Handler:         h.Handler.WithAttrs(attrs),
		additionalAttrs: h.additionalAttrs,
	}
}

// Helper function to create error attributes
func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

// Helper function for structured logging with context
func LogWithContext(ctx context.Context, level slog.Level, msg string, args ...any) {
	// Get trace ID from context if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		args = append(args, slog.String("trace_id", fmt.Sprint(traceID)))
	}

	slog.Log(ctx, level, msg, args...)
}

// Helper function to log errors with stack trace
func LogError(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		args = append(args, ErrAttr(err))
	}
	LogWithContext(ctx, slog.LevelError, msg, args...)
}

// Helper function to create JSON string from struct
func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("error marshaling to JSON: %v", err)
	}
	return string(b)
}

// UnaryClientInterceptor creates a client interceptor for logging
func UnaryClientInterceptor(logger *slog.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		traceID := generateTraceID()

		// Add trace ID to outgoing context
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md.Set("x-trace-id", traceID)
		newCtx := metadata.NewOutgoingContext(ctx, md)

		startTime := time.Now()

		// Create logger with trace ID
		reqLogger := logger.With(
			"trace_id", traceID,
			"method", method,
			"start_time", startTime.Format(time.RFC3339),
		)

		// Log request
		reqLogger.Info("gRPC client request",
			"request", fmt.Sprintf("%+v", req))

		// Make the call
		err := invoker(newCtx, method, req, reply, cc, opts...)

		// Log completion
		duration := time.Since(startTime)
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			}
		}

		reqLogger.Info("gRPC client response",
			"status_code", statusCode.String(),
			"duration_ms", duration.Milliseconds(),
			"error", err,
			"response", fmt.Sprintf("%+v", reply),
		)

		return err
	}
}

func generateTraceID() string {
	return uuid.New().String()
}
