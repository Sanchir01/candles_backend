package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	runtime "github.com/Sanchir01/candles_backend/internal/gql/generated"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

// Candles is the resolver for the candles field.
func (r *mutationResolver) Candles(ctx context.Context) (*model.CandlesMutation, error) {
	return &model.CandlesMutation{}, nil
}

// CandlesMutation returns runtime.CandlesMutationResolver implementation.
func (r *Resolver) CandlesMutation() runtime.CandlesMutationResolver {
	return &candlesMutationResolver{r}
}

type candlesMutationResolver struct{ *Resolver }