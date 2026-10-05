// Command api serves the goland-eda REST and GraphQL APIs.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sathishkottravel/goland-eda/backend/api/internal/config"
	"github.com/sathishkottravel/goland-eda/backend/api/internal/docs"
	"github.com/sathishkottravel/goland-eda/backend/api/internal/graph"
	"github.com/sathishkottravel/goland-eda/backend/api/internal/rest"
	"github.com/sathishkottravel/goland-eda/backend/api/internal/service"
)

func main() {
	cfg := config.Load()

	// `api healthcheck` is used by the container healthcheck (distroless has no curl).
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck(cfg.Port))
	}

	if err := run(cfg); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run(cfg config.Config) error {
	greeter := service.NewGreeter()

	mux := http.NewServeMux()
	rest.New(greeter).Register(mux)
	mux.Handle("POST /graphql", graph.Handler(graph.NewSchema(greeter)))
	docs.Register(mux)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("api listening", "addr", srv.Addr, "kafka_brokers", cfg.KafkaBrokers)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func healthcheck(port string) int {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:" + port + "/health")
	if err != nil {
		return 1
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
