package order

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *Service) CreateOrder(
	ctx context.Context, tx pgx.Tx, userID uuid.UUID, status string, productsId []uuid.UUID, quantity []int, price []int,
) (*uuid.UUID, error) {

	orderID, err := s.repo.CreateOrder(ctx, userID, status, tx)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.CreateOrderItem(ctx, orderID, productsId, quantity, price, tx)
	if err != nil {
		return nil, err
	}

	return &orderID, err
}
