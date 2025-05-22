package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	RepositoryUser *RepositoryUser
}

func NewService(RepositoryUser *RepositoryUser) *Service {
	return &Service{
		RepositoryUser: RepositoryUser,
	}
}

func (s *Service) UserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.RepositoryUser.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *Service) UserByEmail(ctx context.Context, email string, password string) (*model.User, error) {
	usersdb, err := s.RepositoryUser.GetByEmail(ctx, email)
	if err != nil {

		return nil, err
	}
	decodepass, err := base64.StdEncoding.DecodeString(usersdb.Password)
	if err != nil {
		return nil, fmt.Errorf("Неправильный пароль")
	}
	verifypass := VerifyPassword(
		decodepass,
		password,
	)
	if verifypass == false {
		return nil, fmt.Errorf("Неправильный пароль")
	}

	return usersdb, nil
}

func (s *Service) UserByPhone(ctx context.Context, phone string) (*model.User, error) {
	usersdb, err := s.RepositoryUser.GetByPhone(ctx, phone)
	if err != nil {

		return nil, err
	}

	return usersdb, nil
}

func (s *Service) Registrations(ctx context.Context, password, phone, title, email string, tx pgx.Tx) error {
	_, err := s.RepositoryUser.GetByPhone(ctx, phone)
	if err == nil {
		slog.Error("User with this phone already exists")
		return errors.New("user with this phone already exists")
	}

	_, err = s.RepositoryUser.GetByEmail(ctx, email)
	if err == nil {
		slog.Error("User with this email already exists")
		return errors.New("user with this slug already exists")
	}

	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(900000) + 100000

	if err := s.RepositoryUser.SetConfirmationCode(ctx, email, randomNumber); err != nil {
		return err
	}
	if err := SendMail(randomNumber); err != nil {
		return err
	}
	return nil
}
func (s *Service) ConfirmRegister(ctx context.Context, password, phone, title, email, code string, tx pgx.Tx) (*model.User, error) {
	oldcode, err := s.RepositoryUser.GetUserCodeByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if oldcode != code {
		return nil, errors.New("confirmation code is invalid")
	}
	if err := s.RepositoryUser.DeleteUserCodeByEmail(ctx, email); err != nil {
		return nil, err
	}
	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	user, err := s.RepositoryUser.CreateUser(ctx, title, phone, email, "user", hashedPassword, tx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
