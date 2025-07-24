package redirector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/onemoreslacker/shortener/config"
	rpb "github.com/onemoreslacker/shortener/internal/api/proto/redirectorpb"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewHTTPProxy(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	_ *grpc.Server,
	cfg *config.Config,
	log *zap.Logger,
) *http.Server {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(setStatus),
		runtime.WithForwardResponseRewriter(responseEnvelope),
	)

	srv := &http.Server{
		Addr: net.JoinHostPort("", cfg.Serving.RedirectorHTTP),
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
			if err := rpb.RegisterRedirectorHandlerFromEndpoint(
				context.Background(), mux,
				net.JoinHostPort("", cfg.Serving.RedirectorGRPC),
				opts,
			); err != nil {
				return fmt.Errorf("redirector: failed to register redirector handler from endpoint: %w", err)
			}

			go func() {
				log.Info("redirector: starting http proxy server", zap.String("addr", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("redirector: http proxy server error", zap.Error(err))
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

func setStatus(_ context.Context, w http.ResponseWriter, m protoreflect.ProtoMessage) error {
	switch v := m.(type) {
	case *rpb.RedirectResponse:
		w.Header().Set("Location", v.SourceUrl)
		w.WriteHeader(http.StatusMovedPermanently)
	}
	return nil
}

func responseEnvelope(_ context.Context, response proto.Message) (any, error) {
	switch response.(type) {
	case *rpb.RedirectResponse:
		return http.NoBody, nil
	}
	return response, nil
}
