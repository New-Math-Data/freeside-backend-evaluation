-- Document Queries

-- name: CreateDocument :one
INSERT INTO documents (name, current_version)
VALUES ($1, 1)
RETURNING id, name, current_version, created_at;

-- name: GetDocument :one
SELECT id, name, current_version, created_at
FROM documents
WHERE id = $1;

-- name: UpdateDocumentVersion :exec
UPDATE documents
SET current_version = $2
WHERE id = $1;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1;

-- Document Version Queries

-- name: CreateDocumentVersion :one
INSERT INTO document_versions (document_id, version, patch, inverted_patch)
VALUES ($1, $2, $3, $4)
RETURNING document_id, version, patch, inverted_patch, created_at;

-- name: GetDocumentVersion :one
-- TODO: Implement this query
-- Retrieve a specific version record for a document
-- Parameters: document_id (UUID), version (INT)
-- Returns: document_id, version, patch, inverted_patch, created_at
SELECT document_id, version, patch, inverted_patch, created_at
FROM document_versions
WHERE document_id = $1 AND version = $2;

-- name: GetDocumentVersions :many
-- TODO: Implement this query
-- Retrieve all versions for a document, ordered by version ascending
-- This is needed to reconstruct the document at any version
-- Parameters: document_id (UUID)
-- Returns: all columns, ordered by version ASC
SELECT document_id, version, patch, inverted_patch, created_at
FROM document_versions
WHERE document_id = $1
ORDER BY version ASC;

-- name: GetVersionHistory :many
-- TODO: Implement this query
-- Retrieve version metadata (version number and created_at) for displaying history
-- Parameters: document_id (UUID)
-- Returns: version, created_at, ordered by version DESC (newest first)
SELECT version, created_at
FROM document_versions
WHERE document_id = $1
ORDER BY version DESC;

-- name: GetLatestVersion :one
-- TODO: Implement this query
-- Get the latest version number for a document
-- Parameters: document_id (UUID)
-- Returns: the maximum version number
-- Hint: Use MAX() aggregate function
SELECT MAX(version) as version
FROM document_versions
WHERE document_id = $1;
