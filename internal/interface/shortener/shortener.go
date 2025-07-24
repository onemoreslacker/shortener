package shortener

import (
	"context"
	"fmt"

	"github.com/onemoreslacker/shortener/config"
	spb "github.com/onemoreslacker/shortener/internal/api/proto/shortenerpb"
	"github.com/onemoreslacker/shortener/internal/domain/entity"
	"github.com/onemoreslacker/shortener/internal/domain/repository"
	"github.com/onemoreslacker/shortener/pkg/shortid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	links repository.Links
	cfg   *config.Config
	log   *zap.Logger
	spb.UnimplementedShortenerServer
}

func NewServer(
	links repository.Links,
	cfg *config.Config,
	log *zap.Logger,
) *Server {
	return &Server{
		links: links,
		cfg:   cfg,
		log:   log,
	}
}

func (s *Server) ShortenURL(ctx context.Context, req *spb.ShortenRequest) (*spb.ShortenResponse, error) {
	url := req.GetSourceUrl()

	shortID := shortid.Encode(s.cfg.ShorteningPolicy.Length)
	link := entity.NewLink(url, shortID)

	if err := s.links.Add(ctx, link); err != nil {
		s.log.Error("shortener: failed to add shorten link", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to add shorten link")
	}

	shortURL := fmt.Sprintf("%s/%s", s.cfg.ShorteningPolicy.BaseURL, shortID)

	s.log.Info(
		"shortener: add shorten link",
		zap.String("source", link.SourceURL),
		zap.String("short", link.ShortURL),
	)

	return &spb.ShortenResponse{ShortUrl: shortURL}, nil
}
