package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"
	responseErr "github.com/Sanchir01/candles_backend/pkg/lib/api/response"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

// CandleByID is the resolver for the candleById field.
func (r *candlesQueryResolver) CandleByID(ctx context.Context, obj *model.CandlesQuery, input model.CandlesByIDInput) (model.CandlesByIDResult, error) {
	candle, err := r.candlesStr.CandlesById(ctx, input.ID)
	if err != nil {
		return responseErr.NewInternalErrorProblem("не удалось получить свечу"), err
	}
	return model.CandlesByIDOk{Candle: candle}, nil
}