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

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/slemeshaev/todo/internal/handler"
	"github.com/slemeshaev/todo/internal/repository"
	todo "github.com/slemeshaev/todo/internal/server"
	"github.com/slemeshaev/todo/internal/service"
	"github.com/spf13/viper"
)

// @title Todo API
// @version 1.0
// @description API Server for Todo app

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := initConfig(); err != nil {
		slog.Error("failed to init config", "err", err)
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("failed to load env", "err", err)
		os.Exit(1)
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	jwtKey := os.Getenv("JWT_SIGNING_KEY")

	var missing []string

	if dbPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}

	if jwtKey == "" {
		missing = append(missing, "JWT_SIGNING_KEY")
	}

	if len(missing) > 0 {
		slog.Error("missing required env", "vars", missing)
		os.Exit(1)
	}

	// HS256 requires a key of at least 256 bits (RFC 7518, section 3.2).
	const minJWTKeyLen = 32

	if len(jwtKey) < minJWTKeyLen {
		slog.Error("JWT_SIGNING_KEY is too short", "min", minJWTKeyLen, "got", len(jwtKey))
		os.Exit(1)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: dbPassword,
	})
	if err != nil {
		slog.Error("failed to init db", "err", err)
		os.Exit(1)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, jwtKey)
	handlers := handler.NewHandler(services)

	srv := todo.NewServer(viper.GetString("port"), handlers.InitRouters())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	slog.Info("server started", "port", viper.GetString("port"))

	exitCode := 0

	select {
	case <-ctx.Done():
		slog.Info("server shutting down")
	case err := <-errCh:
		exitCode = 1
		slog.Error("http server failed", "err", err)
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown failed", "err", err)
	}

	if err := db.Close(); err != nil {
		slog.Error("db close failed", "err", err)
	}

	slog.Info("server stopped")
	os.Exit(exitCode)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
