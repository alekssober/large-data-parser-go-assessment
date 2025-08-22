package main

import (
	"context"
	"flag"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aleks/products-service/internal/config"
	"github.com/aleks/products-service/internal/csvimporter"
	"github.com/aleks/products-service/internal/db"
	"github.com/aleks/products-service/internal/product"
)

func main() {
	file := flag.String("file", "", "path to CSV file")
	flag.Parse()
	if *file == "" {
		log.Fatal("-file is required")
	}

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
	imp := csvimporter.New(repo)
	if err := imp.FromCSV(ctx, *file); err != nil {
		log.Fatalf("import: %v", err)
	}

	log.Println("import completed ✔")
}
