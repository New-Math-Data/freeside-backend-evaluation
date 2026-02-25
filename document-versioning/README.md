# Document Versioning API - Take Home Assessment

This is a take-home coding assessment for backend engineering candidates.

## Quick Start

1. **Install dependencies:**
   ```bash
   make deps
   ```

2. **Start PostgreSQL:**
   ```bash
   docker compose up -d postgres
   ```

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Generate code:**
   ```bash
   make generate
   ```

5. **Run the server:**
   ```bash
   make run
   ```

## Docker Deployment

To build and run everything in Docker:

```bash
make docker-up
```

The API will be available at http://localhost:8080

## Documentation

See [PROMPT.md](PROMPT.md) for the full assessment instructions.
