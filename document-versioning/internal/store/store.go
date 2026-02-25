package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"document-versioning/internal/database"

	"github.com/google/uuid"
)

// Document represents a document with its current content
type Document struct {
	ID             uuid.UUID
	Name           string
	CurrentVersion int
	Content        map[string]interface{}
	CreatedAt      time.Time
}

// VersionInfo represents metadata about a document version
type VersionInfo struct {
	Version   int
	CreatedAt time.Time
}

// DocumentStore handles document storage and versioning
type DocumentStore struct {
	queries *database.Queries
}

// NewDocumentStore creates a new DocumentStore
func NewDocumentStore(queries *database.Queries) *DocumentStore {
	return &DocumentStore{
		queries: queries,
	}
}

// Create creates a new document with the given content as version 1
// TODO: Implement this method
// 1. Create the document record in the database
// 2. Marshal the content to JSON
// 3. Create a patch from empty {} to the content (this is version 1)
// 4. Store the patch using CreateDocumentVersion
// 5. Return the created document
func (s *DocumentStore) Create(ctx context.Context, name string, content map[string]interface{}) (*Document, error) {
	// Create the document record
	doc, err := s.queries.CreateDocument(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Marshal the content to JSON
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}

	// TODO: Create the initial patch from {} to content
	// Use createPatch() helper function below
	// Store the patch using s.queries.CreateDocumentVersion()

	_ = contentBytes // Use this

	return &Document{
		ID:             doc.ID,
		Name:           doc.Name,
		CurrentVersion: int(doc.CurrentVersion),
		Content:        content,
		CreatedAt:      doc.CreatedAt,
	}, nil
}

// GetCurrent retrieves the current version of a document
// TODO: Implement this method
// 1. Get the document record to find the current version
// 2. Call GetAtVersion with the current version
func (s *DocumentStore) GetCurrent(ctx context.Context, id uuid.UUID) (*Document, error) {
	// Get document metadata
	doc, err := s.queries.GetDocument(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// TODO: Reconstruct the document at current version
	// Call GetAtVersion(ctx, id, int(doc.CurrentVersion))

	_ = doc // Use this

	return nil, fmt.Errorf("not implemented")
}

// GetAtVersion retrieves a document at a specific version
// TODO: Implement this method
// 1. Get all patches up to and including the target version
// 2. Start with an empty document {}
// 3. Apply each patch in order (version 1, 2, ..., N)
// 4. Return the reconstructed document
func (s *DocumentStore) GetAtVersion(ctx context.Context, id uuid.UUID, version int) (*Document, error) {
	// Get document metadata
	doc, err := s.queries.GetDocument(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	if version > int(doc.CurrentVersion) || version < 1 {
		return nil, fmt.Errorf("version %d not found", version)
	}

	// TODO: Get all versions up to the requested version
	// Use s.queries.GetDocumentVersions(ctx, id)
	// Filter to only include versions <= requested version

	// TODO: Reconstruct the document by applying patches
	// Start with base := []byte(`{}`)
	// For each version, apply the patch using applyPatch() helper

	_ = doc // Use this

	return nil, fmt.Errorf("not implemented")
}

// Update updates a document with new content, creating a new version
// TODO: Implement this method
// 1. Get the current document content
// 2. Compute the patch from current to new content
// 3. Store the new version with the patch
// 4. Update the document's current_version
// 5. Return the updated document
func (s *DocumentStore) Update(ctx context.Context, id uuid.UUID, content map[string]interface{}) (*Document, error) {
	// Get current document
	currentDoc, err := s.GetCurrent(ctx, id)
	if err != nil {
		return nil, err
	}

	// Marshal contents
	currentBytes, err := json.Marshal(currentDoc.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal current content: %w", err)
	}

	newBytes, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new content: %w", err)
	}

	// TODO: Create patch from current to new
	// Use createPatch(currentBytes, newBytes)

	// TODO: Create new version record
	// newVersion := currentDoc.CurrentVersion + 1
	// s.queries.CreateDocumentVersion(...)

	// TODO: Update document's current version
	// s.queries.UpdateDocumentVersion(...)

	_ = currentBytes // Use these
	_ = newBytes

	return nil, fmt.Errorf("not implemented")
}

// ListVersions returns the version history for a document
// TODO: Implement this method
// 1. Get all versions for the document
// 2. Return version metadata (version number, created_at)
func (s *DocumentStore) ListVersions(ctx context.Context, id uuid.UUID) (int, []VersionInfo, error) {
	// Verify document exists
	doc, err := s.queries.GetDocument(ctx, id)
	if err != nil {
		return 0, nil, fmt.Errorf("document not found: %w", err)
	}

	// TODO: Get version history using s.queries.GetVersionHistory(ctx, id)
	// Convert to []VersionInfo

	_ = doc // Use this

	return 0, nil, fmt.Errorf("not implemented")
}

// Revert reverts a document to a specific version by creating a new version
// with the content from the target version
// TODO: Implement this method
// 1. Get the document at the target version
// 2. Create a new version with that content (like Update)
func (s *DocumentStore) Revert(ctx context.Context, id uuid.UUID, targetVersion int) (*Document, error) {
	// Get document at target version
	targetDoc, err := s.GetAtVersion(ctx, id, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get target version: %w", err)
	}

	// TODO: Update the document with the target content
	// This creates a new version with the reverted content
	// Use s.Update(ctx, id, targetDoc.Content)

	_ = targetDoc // Use this

	return nil, fmt.Errorf("not implemented")
}

// Helper functions for JSON patch operations

// createPatch creates a JSON patch from base to target, returning both
// the forward patch and inverted patch
func createPatch(base, target []byte) (patch []byte, invertedPatch []byte, err error) {
	// Use jsondiff to create an invertible patch

	return patchBytes, invertedBytes, nil
}

// applyPatch applies a JSON patch to a document
func applyPatch(doc []byte, patchBytes []byte) ([]byte, error) {

	return result, nil
}
