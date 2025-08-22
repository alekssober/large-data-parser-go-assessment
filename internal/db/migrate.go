package db

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var schemaFS embed.FS

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	b, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}
	_, err = pool.Exec(ctx, string(b))
	if err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}
	log.Println("database migrated ✔")
	return nil
}
