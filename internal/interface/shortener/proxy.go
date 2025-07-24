package shortener

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/onemoreslacker/shortener/config"
	spb "github.com/onemoreslacker/shortener/internal/api/proto/shortenerpb"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewHTTPProxy(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	_ *grpc.Server,
	cfg *config.Config,
	log *zap.Logger,
) *http.Server {
	mux := runtime.NewServeMux()

	srv := &http.Server{
		Addr: net.JoinHostPort("", cfg.Serving.ShortenerHTTP),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}
			mux.ServeHTTP(w, r)
		}),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

			if err := spb.RegisterShortenerHandlerFromEndpoint(
				context.Background(), mux,
				net.JoinHostPort("", cfg.Serving.ShortenerGRPC),
				opts,
			); err != nil {
				log.Error("shortener: failed to register shortner handler form endpoint")
				return fmt.Errorf("shortener: failed to register shortener handler from endpoint: %w", err)
			}

			go func() {
				log.Info("shortener: starting http proxy server", zap.String("addr", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("shortener: http proxy server error", zap.Error(err))
					go shutdowner.Shutdown()
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}