package order

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) AllOrders(ctx context.Context) ([]*model.Orders, error) {
	orders, err := s.repo.AllOrders(ctx)
	if err != nil {
		return nil, err
	}
	items, err := utils.MapToGql(orders)
	return items, nil
}
