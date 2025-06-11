package events

import (
	"context"
	"github.com/Sanchir01/candles_backend/pkg/lib/logger/sl"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type EventRepositoryInterface interface {
	CreateEvent(ctx context.Context, eventType, payload string, tx pgx.Tx) (uuid.UUID, error)
	GetManyEvents(ctx context.Context, limit uint64) ([]*EventDB, error)
	SetDone(ctx context.Context, ids []uuid.UUID) error
}
type EventService struct {
	log  *slog.Logger
	repo EventRepositoryInterface
}

func NewEventService(log *slog.Logger, repo EventRepositoryInterface) *EventService {
	return &EventService{
		log:  log,
		repo: repo,
	}
}

func (e *EventService) StartCreateEvent(ctx context.Context, handlePeriod time.Duration, limitEvents uint64) {
	const op = "EventService.StartCreateEvent"

	log := e.log.With(slog.String("op", op))
	ticker := time.NewTicker(handlePeriod)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("stopping event service")
				return

			case <-ticker.C:
				log.Debug("starting process events")

				events, err := e.repo.GetManyEvents(ctx, limitEvents)
				if err != nil {
					log.Error("failed to get new events", sl.Err(err))
					continue
				}

				if len(events) == 0 {
					log.Debug("no events to process")
					continue
				}

				var ids []uuid.UUID
				for _, event := range events {
					ids = append(ids, event.ID)
				}
				e.SendMessage(events[0])
				if err := e.repo.SetDone(ctx, ids); err != nil {
					log.Error("failed to set events done", sl.Err(err))
					continue
				}

				log.Info("successfully processed events", slog.Int("count", len(ids)))
			}
		}
	}()
}
func (e *EventService) SendMessage(event *EventDB) {
	const op = "services.event-sender.SendMessage"

	log := e.log.With(slog.String("op", op))
	log.Info("sending message", slog.Any("event", event))

}
