package color

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/vikstrous/dataloadgen"
	"log"
	"time"
)

type ColorDataLoader struct {
	ColorLoaderByID *dataloadgen.Loader[uuid.UUID, *model.Color]
}

func NewDataLoader(repo *Repository, maxBatch int) *ColorDataLoader {
	loader := dataloadgen.NewLoader(func(ctx context.Context, keys []uuid.UUID) ([]*model.Color, []error) {
		log.Printf("DataLoader keys many load: %v", keys)
		errors := make([]error, len(keys))
		items := make([]*model.Color, len(keys))
		ctxs, cancel := context.WithTimeout(
			context.Background(),
			50*time.Millisecond*time.Duration(len(keys)),
		)
		defer cancel()
		// Получаем цвета по списку ключей
		users, err := repo.ColorByManyId(ctxs, keys)

		if err != nil {
			for i := range errors {
				errors[i] = err
			}
			return nil, errors
		}

		u := make(map[uuid.UUID]*model.Color, len(keys))

		for _, user := range users {
			u[user.ID] = user
		}

		for i, id := range keys {
			items[i] = u[id]
		}

		return items, errors
	}, dataloadgen.WithBatchCapacity(maxBatch), dataloadgen.WithWait(50*time.Millisecond))

	return &ColorDataLoader{
		ColorLoaderByID: loader,
	}
}

func groupColorByID(articles []*model.Color) map[uuid.UUID]*model.Color {
	groups := make(map[uuid.UUID]*model.Color)

	for _, a := range articles {
		groups[a.ID] = a
	}

	return groups
}
