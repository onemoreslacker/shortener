package links

import (
	"context"
	"fmt"

	"github.com/onemoreslacker/shortener/config"
	"github.com/onemoreslacker/shortener/internal/domain/entity"
	"github.com/onemoreslacker/shortener/internal/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Links struct {
	collection *mongo.Collection
}

func NewCollection(client *mongo.Client, cfg *config.Config) *mongo.Collection {
	return client.Database(cfg.Database.Name).Collection(cfg.Collections.Links)
}

func NewRepository(collection *mongo.Collection) repository.Links {
	return &Links{
		collection: collection,
	}
}

func (l *Links) Add(ctx context.Context, link *entity.Link) error {
	if _, err := l.collection.InsertOne(ctx, link); err != nil {
		return fmt.Errorf("links repository: failed to insert link: %w", err)
	}
	return nil
}

func (l *Links) Get(ctx context.Context, shortURL string) (*entity.Link, error) {
	filter := struct {
		ShortURL string `bson:"shortURL"`
	}{
		ShortURL: shortURL,
	}

	link := &entity.Link{}
	if err := l.collection.FindOne(ctx, filter).Decode(link); err != nil {
		return nil, fmt.Errorf("links repository: failed to get link by short url: %w", err)
	}

	return link, nil
}
