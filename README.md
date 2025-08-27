# Products Service (Go + Postgres)

## Limitations & Future Improvements

There are areas that could be enhanced in a production setting:

-   **CSV Importer**: Currently processes line by line with individual upserts. For very large files, this can be slow. Batch inserts or Postgres `COPY` could improve performance.
-   **Validation & Error Reporting**: The importer logs errors but does not provide detailed feedback to users (e.g., failed rows). Schema validation could be extended.
-   **Observability**: No metrics, tracing, or structured logging are included yet. These would be essential for monitoring in production.
-   **Authentication**: The REST API is unauthenticated, which is fine for demo purposes but not secure for real deployments.
-   **Database Migrations**: Schema migrations run inline on startup. A dedicated migration tool should be used in production.
-   **Kubernetes**: Current manifests are minimal (no resource limits, HPA, ingress, or TLS). These would need to be added for real-world usage.

## Quick start (Docker Compose)

```bash
docker compose up --build -d
# import CSV
docker compose exec api /importer -file /data/products.csv
curl http://localhost:8080/products
curl http://localhost:8080/products/summary
```
