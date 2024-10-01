package app

import "github.com/Sanchir01/candles_backend/internal/feature/color"

type Repositories struct {
	ColorRepository *color.Repository
}

func newRepositories(databases *Database) *Repositories {
	return &Repositories{ColorRepository: color.NewRepository(databases.PrimaryDB)}
}
