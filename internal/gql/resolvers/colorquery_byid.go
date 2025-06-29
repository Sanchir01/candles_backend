package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
	responseErr "github.com/Sanchir01/candles_backend/pkg/lib/api/response"
)

// ColorByID is the resolver for the colorById field.
func (r *colorQueryResolver) ColorByID(ctx context.Context, obj *model.ColorQuery, input model.ColorByIDInput) (model.ColorByIDResult, error) {
	color, err := r.env.Services.ColorService.ColorById(ctx, input.ID)
	if err != nil {
		return responseErr.NewInternalErrorProblem("ошибка при получении категории"), err
	}
	return model.ColorByIDOk{
		Colors: color,
	}, nil
}
