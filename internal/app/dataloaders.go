package app

import (
	featurecolor "github.com/Sanchir01/candles_backend/internal/feature/color"
)

const (
	colorLoaderByIDMaxBatch               int = 100
	articleBlockLoaderByArticleIDMaxBatch int = 10
	articleTagLoaderByArticleIDMaxBatch   int = 10
	imageLoaderByIDMaxBatch               int = 100
	projectLoaderByIDMaxBatch             int = 10
	tagLoaderByIDMaxBatch                 int = 100
)

type DataLoaders struct {
	ColorDataLoaderById *featurecolor.ColorDataLoader
}

func newDataLoaders(repos *Repositories) *DataLoaders {
	return &DataLoaders{ColorDataLoaderById: featurecolor.NewDataLoader(repos.ColorRepository, colorLoaderByIDMaxBatch)}
}
