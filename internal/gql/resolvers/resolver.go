package resolver

import (
	"github.com/Sanchir01/candles_backend/internal/app"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	env *app.Env
}

func New(
	env *app.Env,
) *Resolver {
	return &Resolver{
		env: env,
	}
}
