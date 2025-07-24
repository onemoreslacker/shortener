package shortener

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/onemoreslacker/shortener/config"
	spb "github.com/onemoreslacker/shortener/internal/api/proto/shortenerpb"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewGRPCServer(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	srv *Server,
	cfg *config.Config,
	log *zap.Logger,
) *grpc.Server {
	grpcServer := grpc.NewServer()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", net.JoinHostPort("", cfg.Serving.ShortenerGRPC))
			if err != nil {
				return fmt.Errorf("shortener: failed to open tcp connection: %w", err)
			}

			spb.RegisterShortenerServer(grpcServer, srv)

			go func() {
				log.Info("shortener: starting grpc server", zap.String("addr", lis.Addr().String()))
				if err := grpcServer.Serve(lis); !errors.Is(err, grpc.ErrServerStopped) {
					log.Error("shortener: grpc server error", zap.Error(err))
					go shutdowner.Shutdown()
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	return grpcServer
}
