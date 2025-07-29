package metrics

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/onemoreslacker/shortener/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server *http.Server

func New(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	cfg *config.Config,
	log *zap.Logger,
) Server {
	prometheus.MustRegister(RequestsTotal, RequestDuration)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Serving.MetricsHTTP),
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("redirector metrics: starting metrics http server", zap.String("addr", srv.Addr))
				if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
					log.Error("redirector metrics: error in http server", zap.Error(err))
					go shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	},
	)

	return srv
}
