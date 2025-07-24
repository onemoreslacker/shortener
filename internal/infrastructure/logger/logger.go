package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle) *zap.Logger {
	log := zap.Must(zap.NewDevelopment())
	lc.Append(fx.StopHook(log.Sync))
	return log
}
