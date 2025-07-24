package redirector

import (
	"context"
	"errors"

	rpb "github.com/onemoreslacker/shortener/internal/api/proto/redirectorpb"
	"github.com/onemoreslacker/shortener/internal/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	links repository.Links
	log   *zap.Logger
	rpb.UnimplementedRedirectorServer
}

func NewServer(links repository.Links, log *zap.Logger) *Server {
	return &Server{
		links: links,
		log:   log,
	}
}

func (s *Server) RedirectURL(ctx context.Context, req *rpb.RedirectRequest) (*rpb.RedirectResponse, error) {
	shortURL := req.GetShortUrl()

	link, err := s.links.Get(ctx, shortURL)
	if err != nil {
		s.log.Error(
			"redirector: failed to get source url",
			zap.Error(err),
			zap.String("link", shortURL),
		)

		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.Internal, "failed to get source url")

		}

		return nil, status.Errorf(codes.NotFound, "short url not found")
	}

	s.log.Info(
		"redirector: succesfully get source url",
		zap.String("short", link.ShortURL),
		zap.String("source", link.SourceURL),
	)

	return &rpb.RedirectResponse{SourceUrl: link.SourceURL}, nil
}
