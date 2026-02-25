# New Math Data Backend Engineering Take Home Exam

## Introduction

*(This intro is a repeat of the email confirmation instructions. If you have already read through it, feel free to move to the next section.)*

This is a take home exam to focus on coding skills and problem-solving ability using Go. You should have your own Go environment configured for the assignment.

The time allotted for the problem is **8 hours**. You will be responsible for sending a solution prior to this time elapsing.

You can use any available tools at your disposal, including LLMs or AI assistants to help with this problem, as you would use on a normal day-to-day basis as a developer. However, this is a solo exercise, so please do not use a friend.

We don't expect you to be an expert in all of the concepts used here. Part of the exercise is to show your ability to look up and apply common patterns in unknown scenarios.

The problem is judged on the following criteria, in order of importance:

1. Submission of the assignment within the allowed time
2. Working Go code that compiles and runs
3. Correct JSON patch creation and application
4. Version history with bidirectional traversal
5. Proper use of provided tooling (sqlc, OpenAPI types)
6. Successful Docker deployment
7. Code quality, error handling, Go idioms
8. Tests (bonus)

Note that submitting a solution that only meets some of the requirements is preferred to incomplete code.

Hints are provided in the problem statement to aid in solving the problem. Only general familiarity with the Go programming language is assumed for this problem.

---

## Problem

At New Math Data, we work with geospatial data that changes over time. A critical requirement is the ability to track changes to JSON documents, allowing users to view historical states and revert to previous versions when needed.

Your task is to implement a **Document Versioning API** that stores JSON documents with full version history using [RFC 6902 JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902). This is the same approach we use in our production Freeside platform.

### What is JSON Patch?

JSON Patch is a format for describing changes to a JSON document. Instead of storing full copies of each version, we store:

- **Forward patch**: Operations to transform version N → version N+1
- **Inverse patch**: Operations to transform version N+1 → version N (for reverting)

This allows efficient storage while maintaining full bidirectional version history.

### Example

Given an original document:
```json
{"name": "Acme Corp", "employees": 50}
```

And a new version:
```json
{"name": "Acme Corporation", "employees": 75, "active": true}
```

The forward patch would be:
```json
[
  {"op": "replace", "path": "/name", "value": "Acme Corporation"},
  {"op": "replace", "path": "/employees", "value": 75},
  {"op": "add", "path": "/active", "value": true}
]
```

---

## Requirements

### API Endpoints

Implement the following REST API endpoints (OpenAPI spec provided):

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/documents` | Create a new document (version 1) |
| `PUT` | `/api/v1/documents/{id}` | Update document, creates new version |
| `GET` | `/api/v1/documents/{id}` | Get current document state |
| `GET` | `/api/v1/documents/{id}?version={n}` | Get document at specific version |
| `GET` | `/api/v1/documents/{id}/versions` | List all versions with metadata |
| `POST` | `/api/v1/documents/{id}/revert` | Revert to a specific version |

### Your Tasks

1. **Complete the sqlc queries** in `migrations/queries/document_versions.sql` (look for `TODO` comments)

2. **Run code generation**:
   ```bash
   make generate
   ```

3. **Implement versioning logic** in `internal/store/store.go`:
   - Create JSON patches when updating documents
   - Store both forward and inverse patches
   - Reconstruct historical versions by applying patches sequentially

4. **Wire up the handlers** in `internal/handler/documents.go`:
   - Connect the generated OpenAPI types to your store logic

5. **Deploy with Docker**:
   ```bash
   make docker-up
   ```

### Database Schema

The database schema is provided via goose migrations. Two tables are used:

```sql
-- Main document record
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    current_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Version history with patches
CREATE TABLE document_versions (
    document_id UUID REFERENCES documents(id),
    version INT NOT NULL,
    patch JSONB NOT NULL,
    inverted_patch JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (document_id, version)
);
```

### Reconstructing Historical Versions

To get a document at version N:
1. Start with an empty document `{}`
2. Apply patches for versions 1, 2, ..., N in order
3. Return the resulting document

---

## Getting Started

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Make

### Setup

1. **Install development tools**:
   ```bash
   make deps
   ```

2. **Start the database**:
   ```bash
   docker compose up -d postgres
   ```

3. **Run migrations**:
   ```bash
   make migrate-up
   ```

4. **Generate code** (after completing sqlc queries):
   ```bash
   make generate
   ```

5. **Run locally**:
   ```bash
   make run
   ```

### Docker Deployment

To build and run the complete stack:

```bash
make docker-up
```

The API will be available at `http://localhost:8080`.

To stop:
```bash
make docker-down
```

---

## Hints

### Useful Libraries

- **[wI2L/jsondiff](https://github.com/wI2L/jsondiff)** - Create invertible JSON patches
  ```go
  import "github.com/wI2L/jsondiff"
  
  patch, err := jsondiff.CompareJSON(oldDoc, newDoc, jsondiff.Invertible())
  invertedPatch, err := patch.Invert()
  ```

- **[evanphx/json-patch](https://github.com/evanphx/json-patch)** - Apply JSON patches
  ```go
  import jsonpatch "github.com/evanphx/json-patch/v5"
  
  patch, err := jsonpatch.DecodePatch(patchBytes)
  modified, err := patch.Apply(originalBytes)
  ```

### Code Generation

- **sqlc** generates type-safe Go code from SQL queries. After editing `document_versions.sql`, run `make generate`.
- **oapi-codegen** generates Gin server types from OpenAPI. The spec is pre-configured.

### Error Handling

- Return appropriate HTTP status codes (400, 404, 500)
- Version not found should return 404
- Invalid JSON should return 400

---

## How to Submit

Please email your completed source files back to the recruiter in the allotted timeframe. You can ZIP your working directory and send that directly.

**Include:**
- All source files
- Your generated code (from `make generate`)
- Any tests you wrote

**Do NOT include:**
- Binary executables (spam filters will reject)
- The `postgres` data directory

**Please use "Reply All" so the materials go to both the recruiter and the hiring manager.**

---

## Evaluation Checklist

- [ ] Code compiles with `go build ./...`
- [ ] `make generate` succeeds
- [ ] `make docker-up` starts the API successfully
- [ ] Can create a document via POST
- [ ] Can update a document and version increments
- [ ] Can retrieve historical versions
- [ ] Can list version history
- [ ] Can revert to a previous version
- [ ] Code is well-organized and follows Go conventions
