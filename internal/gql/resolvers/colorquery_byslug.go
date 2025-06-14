package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
	responseErr "github.com/Sanchir01/candles_backend/pkg/lib/api/response"
)

// ColorBySlug is the resolver for the colorBySlug field.
func (r *colorQueryResolver) ColorBySlug(ctx context.Context, obj *model.ColorQuery, input model.ColorBySlugInput) (model.ColorBySlugResult, error) {
	colors, err := r.env.Services.ColorService.ColorBySlug(ctx, input.Slug)
	if err != nil {
		return responseErr.NewInternalErrorProblem("error for getting color by slug"), err
	}
	return model.ColorBySlugOk{
		Colors: colors,
	}, nil
}
