package user

import (
	"context"
	"errors"
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
func (s *Service) UserByEmail(ctx context.Context, email string, password string) (*model.User, error) {
	usersdb, err := s.repository.GetByEmail(ctx, email)
	if err != nil {

		return nil, err
	}
	verifypass := VerifyPassword([]byte(usersdb.Password), password)
	if !verifypass {
		return nil, err
	}
	return usersdb, nil
}
func (s *Service) UserByPhone(ctx context.Context, phone string) (*model.User, error) {
	usersdb, err := s.repository.GetByPhone(ctx, phone)
	if err != nil {

		return nil, err
	}

	return usersdb, nil
}

func (s *Service) Registrations(ctx context.Context, password, phone, title, email string, tx pgx.Tx) (*model.User, error) {
	slug, err := utils.Slugify(title)
	if err != nil {
		return nil, err
	}

	_, err = s.repository.GetByPhone(ctx, phone)
	if err == nil {
		slog.Error("User with this phone already exists")
		return nil, errors.New("user with this phone already exists")
	}

	_, err = s.repository.GetByEmail(ctx, email)
	if err == nil {
		slog.Error("User with this slug already exists")
		return nil, errors.New("user with this slug already exists")
	}

	_, err = s.repository.GetBySlug(ctx, slug)
	if err == nil {
		slog.Error("User with this slug already exists")
		return nil, errors.New("user with this slug already exists")
	}

	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		slog.Error("password hash error", err.Error())
		return nil, err
	}

	usersdb, err := s.repository.CreateUser(ctx, title, phone, slug, email, "user", hashedPassword, tx)
	if err != nil {
		return nil, err
	}

	return usersdb, nil
}
