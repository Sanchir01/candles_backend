package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
	responseErr "github.com/Sanchir01/candles_backend/pkg/lib/api/response"
	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
)

// CreateColor is the resolver for the createColor field.
func (r *colorMutationResolver) CreateColor(ctx context.Context, obj *model.ColorMutation, input model.CreateColorInput) (model.ColorCreateResult, error) {
	slug, err := utils.Slugify(input.Title)
	if err != nil {
		return responseErr.NewInternalErrorProblem("error for creating slug"), nil
	}
	id, err := r.colorStr.CreateColor(ctx, input.Title, slug)
	if err != nil {
		r.lg.Error(err.Error())
		return responseErr.NewInternalErrorProblem("error for creating color"), nil
	}
	return model.ColorCreateOk{ID: id}, nil
}