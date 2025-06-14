package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	runtime "github.com/Sanchir01/candles_backend/internal/gql/generated"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

// Color is the resolver for the color field.
func (r *queryResolver) Color(ctx context.Context) (*model.ColorQuery, error) {
	return &model.ColorQuery{}, nil
}

// ColorQuery returns runtime.ColorQueryResolver implementation.
func (r *Resolver) ColorQuery() runtime.ColorQueryResolver { return &colorQueryResolver{r} }

type colorQueryResolver struct{ *Resolver }
