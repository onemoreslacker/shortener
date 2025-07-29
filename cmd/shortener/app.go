package main

import (
	"github.com/onemoreslacker/shortener/config"
	"github.com/onemoreslacker/shortener/internal/infrastructure/logger"
	metrics "github.com/onemoreslacker/shortener/internal/infrastructure/metrics/shortener"
	"github.com/onemoreslacker/shortener/internal/infrastructure/mongo"
	"github.com/onemoreslacker/shortener/internal/infrastructure/persistence/links"
	"github.com/onemoreslacker/shortener/internal/interface/shortener"
	"go.uber.org/fx"
)

func BuildApp() fx.Option {
	return fx.Options(
		logger.FxModule(),
		config.FxModule(),
		metrics.FxModule(),
		mongo.FxModule(),
		links.FxModule(),
		shortener.FxModule(),
	)
}
