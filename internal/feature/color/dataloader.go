package color

import (
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/vikstrous/dataloadgen"
)

type ColorDataLoader struct {
	ColorLoader *dataloadgen.Loader[string, *model.Color]
	ColorRepo   *Repository
}

func NewDataLoader(repo Repository, maxBatch int) *ColorDataLoader {
	//loader := dataloadgen.NewLoader(func(ctx context.Context, keys []string) (ret []*model.Color, errs []error) {
	//
	//	errs = make([]error, len(keys))
	//
	//	users, err := repo.AllColor(ctx)
	//	if err != nil {
	//		return nil, []error{err}
	//	}
	//	ret, err = MapColorToGql(users)
	//	return ret, nil
	//}, dataloadgen.WithBatchCapacity(maxBatch), dataloadgen.WithWait(4*time.Millisecond))
	return &ColorDataLoader{ColorLoader: nil}
}
