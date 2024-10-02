package user

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) UserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UserByPhone(ctx context.Context, phone string) (*model.User, error) {
	usersdb, err := s.repository.GetByPhone(ctx, phone)
	if err != nil {

		return nil, err
	}

	return usersdb, nil
}
