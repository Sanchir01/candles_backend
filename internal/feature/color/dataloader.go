package color

import (
	"context"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/vikstrous/dataloadgen"
	"time"
)

type ColorDataLoader struct {
	ColorLoader *dataloadgen.Loader[string, string]
}
type LoaderByIDRepository interface {
	AllColor(ctx context.Context) ([]model.Color, error)
}

func NewDataLoader(repo LoaderByIDRepository, maxBatch int) *ColorDataLoader {
	loader := dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]string, []error) {
		items := make([]string, len(keys))
		errs := make([]error, len(keys))

		for i, key := range keys {
			items[i] = key
			if key == "errorKey" {
				errs[i] = fmt.Errorf("произошла ошибка с ключом: %s", key)
			} else {
				errs[i] = nil
			}
		}
		return items, errs
	},
		dataloadgen.WithBatchCapacity(maxBatch),
		dataloadgen.WithWait(3*time.Millisecond),
	)
	return &ColorDataLoader{ColorLoader: loader}
}

func (d *ColorDataLoader) GetColor(ctx context.Context, key string) (string, error) {
	return d.ColorLoader.Load(ctx, key)
}

func (d *ColorDataLoader) GetColors(ctx context.Context, keys []string) ([]string, []error) {
	allColor, err := d.ColorLoader.LoadAll(ctx, keys)
	if err != nil {
		return nil, []error{err}
	}
	return allColor, nil
}
