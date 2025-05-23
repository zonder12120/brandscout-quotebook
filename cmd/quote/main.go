package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zonder12120/brandscout-quotebook/internal/config"
	"github.com/zonder12120/brandscout-quotebook/internal/rest"
	"github.com/zonder12120/brandscout-quotebook/internal/rest/handler"
	"github.com/zonder12120/brandscout-quotebook/internal/service"
	"github.com/zonder12120/brandscout-quotebook/internal/storage"
	"github.com/zonder12120/brandscout-quotebook/pkg/logger"
)

const gracefulShutdownTimeout = 5 * time.Second

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.LogLevel)

	quoteStorage := storage.NewInMemory(cfg.QuotesLimit)
	quoteService := service.NewQuoteService(quoteStorage)
	quoteHandler := handler.New(quoteService, log)

	router := rest.NewRouter(quoteHandler, log)

	addr := cfg.Port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Msgf("Starting server on %s", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Shutting down server ...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("Server stopped gracefully")
	}
}
