package featurecolor

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
)

type LoaderByIDRepository interface {
	Find(ctx context.Context) ([]model.Color, error)
}

type LoaderByIDKey struct {
	ID uuid.UUID
}

func NewConfiguredLoaderById() {}
