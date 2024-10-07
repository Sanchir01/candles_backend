package candles

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repository *Repository
	storages   *Storage
}

func NewService(repository *Repository, storages *Storage) *Service {
	return &Service{
		repository,
		storages,
	}
}
func (s *Service) AllCandles(ctx context.Context, sort *model.CandlesSortEnum) ([]*model.Candles, error) {

	candles, err := s.repository.AllCandles(ctx, sort)
	gqlCandles, err := MapCandlesToGql(candles)

	return gqlCandles, err
}

func (s *Service) CreateCandles(ctx context.Context, categoryID, colorID uuid.UUID, title string, images []*graphql.Upload, price int) (uuid.UUID, error) {
	conn, err := s.repository.primaryDB.Acquire(ctx)
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
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()
	slug, err := utils.Slugify(title)
	if err != nil {
		return uuid.Nil, err
	}
	_, err = s.repository.CandlesBySlug(ctx, slug)
	if err == nil {
		return uuid.Nil, err
	}
	imagesUrl, err := s.storages.PutObjects(ctx, "candles", images)
	if err != nil {
		return uuid.Nil, nil
	}
	id, err := s.repository.CreateCandles(ctx, categoryID, colorID, title, slug, imagesUrl, price, tx)
	if err != nil {
		s.storages.DeleteObjects(ctx, "candles", images)
		return uuid.Nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (s *Service) CandlesById(ctx context.Context, id uuid.UUID) (*model.Candles, error) {

	candles, err := s.repository.CandlesById(ctx, id)

	return candles, err
}

func (s *Service) CandlesBySlug(ctx context.Context, title string) (*model.Candles, error) {
	slug, err := utils.Slugify(title)
	if err != nil {
		return nil, err
	}
	candles, err := s.repository.CandlesBySlug(ctx, slug)

	return candles, err
}
