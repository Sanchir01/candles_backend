package order

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/feature/events"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo       *Repository
	eventsrepo *events.Repository
}

func NewService(repo *Repository, eventsrepo *events.Repository) *Service {
	return &Service{
		repo, eventsrepo,
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

func (s *Service) AllUserOrders(ctx context.Context, id uuid.UUID) ([]*model.Orders, error) {
	orders, err := s.repo.OrdersByUserId(ctx, id)
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
	_, err = s.eventsrepo.CreateEvent(ctx, utils.EventSavedOrder, orderID.String(), tx)
	if err != nil {
		return nil, err
	}
	uuids, err := s.repo.CreateOrderItem(ctx, orderID, productsId, quantity, price, tx)
	if err != nil {
		return nil, err
	}

	return uuids, err
}

func (s *Service) OrderById(ctx context.Context, id uuid.UUID) (string, error) {
	status, err := s.repo.GetOrderById(ctx, id)
	if err != nil {
		return "", err
	}
	return status, nil
}
