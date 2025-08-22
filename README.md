# Products Service (Go + Postgres)

## Quick start (Docker Compose)

```bash
docker compose up --build -d
# import CSV
docker compose exec api /importer -file /data/products.csv
curl http://localhost:8080/products
curl http://localhost:8080/products/summary
```
