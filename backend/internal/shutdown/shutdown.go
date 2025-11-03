// internal/shutdown/shutdown.go
package shutdown

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Stoppable interface {
	Stop(ctx context.Context) error
}

// Wait blocks until SIGINT/SIGTERM and then gracefully stops the app.
func Wait(app Stoppable, log *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stop
	log.Info("received shutdown signal", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Error("graceful shutdown failed", slog.Any("err", err))
	} else {
		log.Info("server stopped gracefully")
	}
}
