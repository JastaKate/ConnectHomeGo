package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectHome/internal/clients/postgres"
	"connectHome/internal/controller"
	"connectHome/internal/repository"
	"connectHome/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	cfg, err := postgres.LoadFromEnv()
	if err != nil {
		return err
	}

	pool, err := postgres.New(context.Background(), cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	routes := createRoutes(pool)
	return startServer(routes)
}

func createRoutes(pool *pgxpool.Pool) *http.ServeMux {
	homeRepository := repository.NewHomeRepository(pool)
	homeService := service.NewHomeService(homeRepository)
	homeController := controller.NewHomeController(homeService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", controller.Health)
	mux.HandleFunc("POST /v2/api/homes", homeController.Create)
	mux.HandleFunc("GET /v2/api/homes", homeController.List)
	return mux
}

func startServer(routes http.Handler) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Server listening on http://localhost:8080")
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			return err
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return err
		}
	}
	return nil
}
