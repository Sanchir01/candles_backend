package category

import (
	"context"
	"fmt"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
)

type Service struct {
	repository CategoryService
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=CategoryService
type CategoryService interface {
	CategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	CategoryBySlug(ctx context.Context, slug string) (*model.Category, error)
	AllCategories(ctx context.Context) ([]model.Category, error)
	CreateCategory(ctx context.Context, title, slug string) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name, slug string) (uuid.UUID, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) AllCategory(ctx context.Context) ([]*model.Category, error) {
	allCategory, err := s.repository.AllCategories(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := utils.MapToGql(allCategory)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) CategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	return s.repository.CategoryById(ctx, id)
}

func (s *Service) CategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	return s.repository.CategoryBySlug(ctx, slug)
}

func (s *Service) CreateCategory(ctx context.Context, title string) (uuid.UUID, error) {
	slug, err := utils.Slugify(title)
	if err != nil {
		return uuid.Nil, err
	}
	isExistCategory, err := s.repository.CategoryBySlug(ctx, slug)
	if err == nil {
		return uuid.Nil, fmt.Errorf("категория с slug: %s уже существует", isExistCategory.Slug)
	}
	id, err := s.repository.CreateCategory(ctx, title, slug)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Service) DeleteCategoryById(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	id, err := s.repository.DeleteCategory(ctx, id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
