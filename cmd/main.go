package main

import (
	"context"
	"log/slog"
	"merch-shop/internal/api"
	"merch-shop/internal/auth"
	"merch-shop/internal/config"
	"merch-shop/internal/repository"
	"merch-shop/internal/repository/db"
	"merch-shop/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.Info("Start server")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("config.LoadConfig", "error", err)
		return
	}

	privateKey, err := config.ParsePrivateKey(cfg.PrivateKey)
	if err != nil {
		slog.Error("config.ParsePrivateKey", "error", err)
		return
	}

	publicKey, err := config.ParsePublicKey(cfg.PublicKey)
	if err != nil {
		slog.Error("config.ParsePublicKey", "error", err)
		return
	}

	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("db.Connect", "error", err)
		return
	}
	defer db.Close()

	repo := repository.New(db)
	auth := authorization.New(privateKey, publicKey)

	useCase := usecase.New(auth, repo)

	handler := api.NewHTTPHandler(useCase)
	router, err := api.NewRouter(handler, publicKey)
	if err != nil {
		slog.Error("api.NewRouter", "error", err)
		return
	}

	srv := api.NewServer(cfg.ServerPort, router)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//Запускаем сервер
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Service running", "error", err)
			return
		}
	}()

	<-ctx.Done()
	slog.Info("Server stop start")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return
	}

	slog.Info("Server stopped")
}
