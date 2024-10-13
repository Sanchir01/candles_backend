package order

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/bot"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo *Repository
	bot  *bot.Bot
}

func NewService(repo *Repository, bot *bot.Bot) *Service {
	return &Service{
		repo,
		bot,
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
) ([]uuid.UUID, error) {
	totalAmount := 0
	for _, q := range quantity {
		totalAmount += q
	}
	orderID, err := s.repo.CreateOrder(ctx, userID, status, totalAmount, tx)
	if err != nil {
		return nil, err
	}

	uuids, err := s.repo.CreateOrderItem(ctx, orderID, productsId, quantity, price, tx)
	if err != nil {
		return nil, err
	}

	return uuids, err
}
