package common

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
}

var (
	cRED    = color.New(color.FgRed).Add(color.Bold).SprintFunc()
	cGREEN  = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	cYELLOW = color.New(color.FgYellow).Add(color.Bold).SprintFunc()
)

type customHandler struct {
	level slog.Level
}

func (h *customHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := time.Now().Format(time.RFC3339)
	levelStr := r.Level.String()

	// Aplicar colores según el nivel
	switch r.Level {
	case slog.LevelInfo:
		levelStr = cGREEN(levelStr)
	case slog.LevelWarn:
		levelStr = cYELLOW(levelStr)
	case slog.LevelError:
		levelStr = cRED(levelStr)
	}

	// Variables para name y tenant
	var tenantID, tenantName string
	var attrsMsg string

	// Iterar sobre atributos
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "tenant_id":
			tenantID = fmt.Sprintf("%v", a.Value)
		case "tenant_name":
			tenantName = fmt.Sprintf("%v", a.Value)
		default:
			attrsMsg += fmt.Sprintf(" %s=%v", a.Key, a.Value)
		}
		return true
	})

	// Construcción del mensaje final
	// msg := fmt.Sprintf("%[2]s=[%[1]s] level=%[2]s", timestamp, levelStr)
	msg := fmt.Sprintf("%s=[%s]", levelStr, timestamp)
	if tenantID != "" {
		msg += fmt.Sprintf(" tenant_id=%s", cGREEN(tenantID))
	}

	if tenantName != "" {
		msg += fmt.Sprintf(" tenant_name=%s", cGREEN(tenantName))
	}

	msg += fmt.Sprintf(" msg=%q%s", r.Message, attrsMsg)
	fmt.Println(msg)
	return nil
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	return h
}

type logger struct {
	logger        *slog.Logger
	tenantManager *TenantConnectionManager
}

func NewLoggerWithTenantManager(tenantManager *TenantConnectionManager) Logger {
	handler := &customHandler{level: slog.LevelInfo}
	l := &logger{
		logger:        slog.New(handler),
		tenantManager: tenantManager,
	}
	return l
}

func NewLogger() Logger {
	return NewLoggerWithTenantManager(nil)
}

func (l *logger) Info(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, l.tenantInfo(ctx)...)
	l.logger.Info(msg, args...)
}

func (l *logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, l.tenantInfo(ctx)...)
	l.logger.Warn(msg, args...)
}

func (l *logger) Error(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, l.tenantInfo(ctx)...)
	l.logger.Error(msg, args...)
}

func (l *logger) tenantInfo(ctx context.Context) []interface{} {
	var args []interface{}
	if l.tenantManager == nil {
		return args
	}
	if tenantDetail, ok := ctx.Value(l.tenantManager.TenantKey).(uuid.UUID); ok {
		tenant, err := l.tenantManager.GetTenantConfig(tenantDetail)
		if err != nil {
			l.Warn(ctx, "Tenant not found")
			return args
		}
		tenantID := tenant.TenantID.String()
		name := tenant.Name
		args = append(args, slog.String("tenant_name", name))
		args = append(args, slog.String("tenant_id", tenantID))
	}
	return args
}

var _ Logger = (*logger)(nil)
