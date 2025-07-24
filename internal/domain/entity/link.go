package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SourceURL string             `bson:"sourceURL"`
	ShortURL  string             `bson:"shortURL"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func NewLink(sourceURL, shortURL string) *Link {
	return &Link{
		SourceURL: sourceURL,
		ShortURL:  shortURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
