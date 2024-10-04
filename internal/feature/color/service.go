package color

import (
	"context"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) AllColor(ctx context.Context) ([]*model.Color, error) {
	colors, err := s.repository.AllColor(ctx)
	if err != nil {
		return nil, err
	}
	gqlcolors, err := MapColorToGql(colors)
	if err != nil {
		return nil, err
	}
	return gqlcolors, nil
}

func (s *Service) CreateColor(ctx context.Context, title string) (uuid.UUID, error) {
	slug, err := utils.Slugify(title)
	if err != nil {
		return uuid.Nil, err
	}
	isExistColor, err := s.repository.ColorBySlug(ctx, slug)
	if err == nil {
		return uuid.Nil, fmt.Errorf("цвет с slug: %s уже существует", isExistColor.Slug)
	}
	id, err := s.repository.CreateColor(ctx, title, slug)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
