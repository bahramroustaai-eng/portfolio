package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"portfolio/internal/debt"
	"portfolio/internal/platform/config"
	"portfolio/internal/platform/postgres"
	"portfolio/internal/portfo"
	transport "portfolio/internal/transport/http"
	"portfolio/internal/user"
	"syscall"
	"time"

	_ "portfolio/docs"
)

// @title 			Portfolio
// @version         1.0
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := user.NewUserRepository(pool)
	userSvc := user.NewUserService(userRepo)
	userHandler := transport.NewUserHandler(userSvc)

	debtRepo := debt.NewRepository(pool)
	debtSvc := debt.NewDebtService(debtRepo, userRepo)
	debtHandler := transport.NewDebtHandler(debtSvc)

	portfoRepo := portfo.NewRepository(pool)
	portfoSvc := portfo.NewService(portfoRepo)
	portfolioHandler := transport.NewPortfolioHandler(portfoSvc)

	router := transport.NewRouter(userHandler, debtHandler, portfolioHandler)

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server running...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
