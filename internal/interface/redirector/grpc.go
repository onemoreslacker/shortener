package redirector

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/onemoreslacker/shortener/config"
	rpb "github.com/onemoreslacker/shortener/internal/api/proto/redirectorpb"
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
			lis, err := net.Listen("tcp", net.JoinHostPort("", cfg.Serving.RedirectorGRPC))
			if err != nil {
				return fmt.Errorf("redirector: failed to open connection: %w", err)
			}

			rpb.RegisterRedirectorServer(grpcServer, srv)

			go func() {
				log.Info("redirector: starting grpc server", zap.String("addr", lis.Addr().String()))
				if err := grpcServer.Serve(lis); !errors.Is(err, grpc.ErrServerStopped) {
					log.Error("redirector: grpc server error", zap.Error(err))
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
