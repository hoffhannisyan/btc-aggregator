package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hoffhannisyan/btc-aggregator/internal/aggregator"
	"github.com/hoffhannisyan/btc-aggregator/internal/api"
	"github.com/hoffhannisyan/btc-aggregator/internal/config"
	"github.com/hoffhannisyan/btc-aggregator/internal/metrics"
	"github.com/hoffhannisyan/btc-aggregator/internal/poller"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("config loaded",
		"poll_interval", cfg.PollInterval,
		"request_timeout", cfg.RequestTimeout,
		"max_retries", cfg.MaxRetries,
		"http_port", cfg.HTTPPort,
	)

	m := metrics.New()
	agg := aggregator.New()
	p := poller.New(cfg, agg, m, logger)
	handler := api.NewHandler(agg)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go p.Start(ctx)

	go func() {
		logger.Info("HTTP server started", "port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}

	logger.Info("server stopped")
}