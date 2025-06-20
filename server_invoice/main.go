package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/server_invoice/api"
	"github.com/LamThanhNguyen/yoyo-store-backend/server_invoice/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	config, err := util.LoadConfig(ctx, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	log.Info().Interface("config", config).Msg("loaded config")

	waitGroup, ctx := errgroup.WithContext(ctx)

	runServer(ctx, waitGroup, config)

	if err = waitGroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("err from wait group")
	}

	log.Info().Msg("application shutdown complete")
}

func runServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
) {
	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	server.SetupRouter() // initialize routes

	server.CreateDirIfNotExist("./invoices")

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:              config.InvoicePort,
		Handler:           server.Router(), // use the Gin router
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// Start HTTP server in goroutine
	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP server at %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("HTTP server failed to serve")
			return err
		}
		return nil
	})

	// Graceful shutdown on context cancel
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP server")
			return err
		}

		log.Info().Msg("HTTP server is stopped")
		return nil
	})
}
