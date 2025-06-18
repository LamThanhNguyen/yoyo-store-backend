package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/LamThanhNguyen/yoyo-store-backend/db/sqlc"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
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

	config, err := LoadConfig(ctx, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	runtimeCfg, err := NewRuntimeConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid config values")
	}

	log.Info().Interface("config", runtimeCfg).Msg("loaded config")

	connPool, err := pgxpool.New(ctx, runtimeCfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	// Run migration to database
	runDBMigration(runtimeCfg.MigrationURL, runtimeCfg.DBSource)

	store := db.NewStore(connPool)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runServer(ctx, waitGroup, runtimeCfg, store)

	if err = waitGroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("err from wait group")
	}

	log.Info().Msg("application shutdown complete")
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config RuntimeConfig,
	store db.Store,
) {
	server, err := NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	server.SetupRouter() // initialize routes

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:              config.ApiPort,
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
