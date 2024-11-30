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

type ServiceCandles struct {
	repository *RepositoryCandles
	storages   *Storage
}

func NewServiceCandles(repository *RepositoryCandles, storages *Storage) *ServiceCandles {
	return &ServiceCandles{
		repository,
		storages,
	}
}
func (s *ServiceCandles) AllCandles(ctx context.Context, sort *model.CandlesSortEnum, filter *model.CandlesFilterInput, pageSize uint, pageNumber uint) ([]*model.Candles, error) {

	candles, err := s.repository.AllCandles(ctx, sort, filter, pageSize, pageNumber)
	if err != nil {
		return nil, err
	}

	gqlCandles, err := MapCandlesToGql(candles)
	if err != nil {
		return nil, err
	}

	return gqlCandles, nil
}
func (s *ServiceCandles) GetTotalCountCandles(ctx context.Context, filter *model.CandlesFilterInput) (uint, error) {
	totalcount, err := s.repository.CountCandles(ctx, filter)
	if err != nil {
		return 0, err
	}
	return totalcount, nil
}
func (s *ServiceCandles) CreateCandles(ctx context.Context, categoryID, colorID uuid.UUID, title, description string, images []*graphql.Upload, price, weight int) (uuid.UUID, error) {
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
	id, err := s.repository.CreateCandles(ctx, categoryID, colorID, title, slug, description, imagesUrl, weight, price, tx)
	if err != nil {
		s.storages.DeleteObjects(ctx, "candles", images)
		return uuid.Nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *ServiceCandles) CandlesById(ctx context.Context, id uuid.UUID) (*model.Candles, error) {
	candles, err := s.repository.CandlesById(ctx, id)
	if err != nil {
		return nil, err
	}
	return candles, err
}

func (s *ServiceCandles) CandlesBySlug(ctx context.Context, title string) (*model.Candles, error) {
	slug, err := utils.Slugify(title)
	if err != nil {
		return nil, err
	}
	candles, err := s.repository.CandlesBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return candles, err
}
