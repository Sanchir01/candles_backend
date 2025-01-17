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
	gqlcolors, err := utils.MapToGql(colors)
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

func (s *Service) ColorBySlug(ctx context.Context, slug string) (*model.Color, error) {
	color, err := s.repository.ColorBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return color, nil
}

func (s *Service) ColorById(ctx context.Context, id uuid.UUID) (*model.Color, error) {
	color, err := s.repository.ColorById(ctx, id)
	if err != nil {
		return nil, err
	}
	return color, nil
}

func (s *Service) UpdateColorById(ctx context.Context, id uuid.UUID, title, slug string) (uuid.UUID, error) {
	colorId, err := s.repository.UpdateCategory(ctx, id, title, slug)
	if err != nil {
		return uuid.Nil, err
	}
	return colorId, nil
}

func (s *Service) DeleteColorById(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	id, err := s.repository.DeleteColor(ctx, id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
