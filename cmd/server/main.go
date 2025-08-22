package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleks/products-service/internal/config"
	"github.com/aleks/products-service/internal/db"
	"github.com/aleks/products-service/internal/httpx"
	"github.com/aleks/products-service/internal/product"
)

func main() {
	cfg := config.FromEnv()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := product.NewPGRepo(pool)
	svc := product.NewService(repo)
	h := product.NewHandler(svc)

	srv := httpx.NewServer(func(r chi.Router) {
		h.Routes(r)
	})

	srv.Addr = ":" + cfg.Port

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShutdown)
	log.Println("server stopped")
}
