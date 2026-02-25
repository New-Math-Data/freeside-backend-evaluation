-- +goose Up
-- +goose StatementBegin

-- Documents table stores the main document record
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    current_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Document versions table stores patches for each version
-- To reconstruct version N, apply patches 1..N to an empty document {}
CREATE TABLE document_versions (
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version INT NOT NULL,
    patch JSONB NOT NULL,           -- Forward patch: transforms version N-1 to N
    inverted_patch JSONB NOT NULL,  -- Inverse patch: transforms version N to N-1
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (document_id, version)
);

-- Index for efficient version lookups
CREATE INDEX idx_document_versions_document_id ON document_versions(document_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;

-- +goose StatementEnd
