package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	runtime "github.com/Sanchir01/candles_backend/internal/gql/generated"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context) (*model.OrderQuery, error) {
	return &model.OrderQuery{}, nil
}

// OrderQuery returns runtime.OrderQueryResolver implementation.
func (r *Resolver) OrderQuery() runtime.OrderQueryResolver { return &orderQueryResolver{r} }

type orderQueryResolver struct{ *Resolver }
