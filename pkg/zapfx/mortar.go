package zapfx

import (
	"context"
	"fmt"

	mortarLog "github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NewMortarLogParams struct {
	fx.In

	Logger *zap.Logger
}

type NewMortarLogResult struct {
	fx.Out

	Logger mortarLog.Logger
}

type mortarProxy struct {
	zap *zap.Logger
}

func NewMortarProxy(p NewMortarLogParams) (NewMortarLogResult, error) {
	return NewMortarLogResult{Logger: &mortarProxy{zap: p.Logger}}, nil
}

func (m *mortarProxy) Trace(ctx context.Context, format string, args ...interface{}) {
	m.zap.Debug(fmt.Sprintf(format, args...))
}

func (m *mortarProxy) Debug(ctx context.Context, format string, args ...interface{}) {
	m.zap.Debug(fmt.Sprintf(format, args...))
}

func (m *mortarProxy) Info(ctx context.Context, format string, args ...interface{}) {
	m.zap.Info(fmt.Sprintf(format, args...))
}

func (m *mortarProxy) Warn(ctx context.Context, format string, args ...interface{}) {
	m.zap.Warn(fmt.Sprintf(format, args...))
}

func (m *mortarProxy) Error(ctx context.Context, format string, args ...interface{}) {
	m.zap.Error(fmt.Sprintf(format, args...))
}

func (m *mortarProxy) Custom(ctx context.Context, level mortarLog.Level, skipAdditionalFrames int, format string, args ...interface{}) {
	m.zap.Warn("warning: level discarded in mortarProxy Custom Logger")
	m.zap.Info(fmt.Sprintf(format, args...))
}

func (m *mortarProxy) WithError(err error) mortarLog.Fields {
	m.zap.Warn("warning: mortarProxy Custom Logger: WithError unimplemented")
	return nil
}

func (m *mortarProxy) WithField(name string, value interface{}) mortarLog.Fields {
	m.zap.Warn("warning: mortarProxy Custom Logger: WithField unimplemented")
	return nil
}

func (m *mortarProxy) Configuration() mortarLog.LoggerConfiguration {
	m.zap.Warn("warning: mortarProxy Custom Logger: Configuration unimplemented")
	return nil
}
