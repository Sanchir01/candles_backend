package suite

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/config"
	"testing"
)

type Suite struct {
	*testing.T
	cfg *config.Config
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.InitConfig()

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})
	ctx.Done()
	return nil, nil
}
