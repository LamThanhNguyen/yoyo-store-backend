package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/frontend/handler"
	"github.com/LamThanhNguyen/yoyo-store-backend/frontend/util"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var Session *scs.SessionManager

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	gob.Register(handler.TransactionData{})

	config, err := util.LoadConfig(ctx, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	runtimeCfg, err := util.NewRuntimeConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid config values")
	}

	log.Info().Interface("config", runtimeCfg).Msg("loaded config")

	// connPool, err := pgxpool.New(ctx, runtimeCfg.DBSource)
	connPool, err := sql.Open("pgx", runtimeCfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	Session = scs.New()
	Session.Lifetime = 24 * time.Hour
	Session.Store = postgresstore.New(connPool)

	tc := make(map[string]*template.Template)

	db_model := models.DBModel{DB: connPool}

	waitGroup, ctx := errgroup.WithContext(ctx)

	runServer(ctx, waitGroup, runtimeCfg, tc, db_model, Session, connPool)

	if err = waitGroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("err from wait group")
	}

	log.Info().Msg("application shutdown complete")
}

func runServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.RuntimeConfig,
	templateCache map[string]*template.Template,
	db models.DBModel,
	session *scs.SessionManager,
	dbConn *sql.DB,
) {
	server, err := handler.NewServer(config, templateCache, db, session)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	server.SetupRouter() // initialize routes

	waitGroup.Go(func() error {
		server.ListenToWsChannel(ctx)
		return nil
	})

	httpServer := &http.Server{
		Addr:              config.FrontendPort,
		Handler:           server.Router(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
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

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Closing DB connection")
		if err := dbConn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close DB connection")
			return err
		}
		return nil
	})
}
