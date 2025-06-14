package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	runtime "github.com/Sanchir01/candles_backend/internal/gql/generated"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

// Category is the resolver for the category.sql field.
func (r *queryResolver) Category(ctx context.Context) (*model.CategoryQuery, error) {
	return &model.CategoryQuery{}, nil
}

// CategoryQuery returns runtime.CategoryQueryResolver implementation.
func (r *Resolver) CategoryQuery() runtime.CategoryQueryResolver { return &categoryQueryResolver{r} }

type categoryQueryResolver struct{ *Resolver }
