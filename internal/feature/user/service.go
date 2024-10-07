package user

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
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

func (s *Service) Registrations(ctx context.Context, title, phone, role string, tx pgx.Tx) (*model.User, error) {

	slug, err := utils.Slugify(title)
	if err != nil {
		return nil, err
	}
	existUser, err := s.UserByPhone(ctx, phone)

	if err == nil {
		slog.Error("User with this phone already exists", existUser)
		return nil, err
	}
	usersdb, err := s.repository.CreateUser(ctx, title, phone, slug, role, tx)
	if err != nil {
		return nil, err
	}

	return usersdb, nil
}
