package color

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/feature/events"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Service struct {
	repository *Repository
	eventrepo  *events.Repository
	primaryDB  *pgxpool.Pool
}

func NewService(repository *Repository, eventrepo *events.Repository, primaryDB *pgxpool.Pool) *Service {
	return &Service{
		repository, eventrepo, primaryDB,
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

func (s *Service) CreateColor(ctx context.Context, title string) (id uuid.UUID, err error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
	}()

	slug, err := utils.Slugify(title)
	if err != nil {
		return uuid.Nil, err
	}
	isExistColor, err := s.repository.ColorBySlug(ctx, slug)
	if err == nil {
		return uuid.Nil, fmt.Errorf("цвет с slug: %s уже существует", isExistColor.Slug)
	}
	id, err = s.repository.CreateColor(ctx, title, slug, tx)
	if err != nil {
		return uuid.Nil, err
	}
	colordata, err := json.Marshal(id)
	if err != nil {
		return uuid.Nil, err
	}
	_, err = s.eventrepo.CreateEvent(ctx, utils.EventSavedColor, string(colordata), tx)
	if err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		slog.Error("failed to commit transaction", "error", err)
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
