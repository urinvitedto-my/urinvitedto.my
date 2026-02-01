package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urinvitedto-my/backend/internal/config"
	"github.com/urinvitedto-my/backend/internal/database"
	"github.com/urinvitedto-my/backend/internal/router"
)

func main() {
	slog.Info("Loading config...")
	cfg, cfgErr := config.LoadCfg()
	if cfgErr != nil {
		slog.Error("Error loading config", "error", cfgErr)
		os.Exit(1)
	}

	slog.Info("Setting up database...")
	db, dbErr := database.Connect(context.Background(), cfg.DatabaseURL)
	if dbErr != nil {
		slog.Error("Database connection failed", "error", dbErr)
		os.Exit(1)
	} else {
		slog.Info("Database connected")
		defer db.Close()
	}

	slog.Info("Setting up routers...")
	r := router.NewRouter(cfg, db)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r.SetupRouter(),
	}

	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		slog.Info("Test server heartbeat", "url", "http://localhost:8080/ping")
		slog.Info("Check server status", "url", "http://localhost:8080/api/v1/")
		slog.Info("Go to this url to test server", "url", "http://localhost:8080/")

		if srvErr := srv.ListenAndServe(); srvErr != nil &&
			srvErr != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", srvErr)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if srvErr := srv.Shutdown(shutdownCtx); srvErr != nil {
		slog.Error("Server forced to shutdown", "error", srvErr)
	} else {
		slog.Info("Server exited properly")
	}
}
