package app

const (
	articleLoaderByIDMaxBatch             int = 100
	articleBlockLoaderByArticleIDMaxBatch int = 10
	articleTagLoaderByArticleIDMaxBatch   int = 10
	imageLoaderByIDMaxBatch               int = 100
	projectLoaderByIDMaxBatch             int = 10
	tagLoaderByIDMaxBatch                 int = 100
)

type DataLoaders struct {
}

func newDataLoaders() {
	return
}
