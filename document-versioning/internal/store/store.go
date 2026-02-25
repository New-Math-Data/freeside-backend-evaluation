package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"document-versioning/internal/database"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wI2L/jsondiff"
)

func convertUUID(source pgtype.UUID) (uuid.UUID, error) {
	if !source.Valid {
		return uuid.Nil, fmt.Errorf("invalid pgtype.UUID")
	}
	result, err := uuid.FromBytes(source.Bytes[:])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return result, nil
}

func setUUID(dest *pgtype.UUID, value uuid.UUID) {
	dest.Bytes = value
	dest.Valid = true
}

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

	// Store the initial doc version
	patch, invertedPatch, err := createPatch([]byte("{}"), contentBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create patch: %w", err)
	}
	_, err = s.queries.CreateDocumentVersion(ctx, database.CreateDocumentVersionParams{
		DocumentID:    doc.ID,
		Version:       1,
		Patch:         patch,
		InvertedPatch: invertedPatch,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document version: %w", err)
	}

	id, err := convertUUID(doc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document ID: %w", err)
	}

	return &Document{
		ID:             id,
		Name:           doc.Name,
		CurrentVersion: int(doc.CurrentVersion),
		Content:        content,
		CreatedAt:      doc.CreatedAt.Time,
	}, nil
}

// GetCurrent retrieves the current version of a document
// 1. Get the document record to find the current version
// 2. Call GetAtVersion with the current version
func (s *DocumentStore) GetCurrent(ctx context.Context, id uuid.UUID) (*Document, error) {
	var docID pgtype.UUID
	setUUID(&docID, id)

	doc, err := s.queries.GetDocument(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	return s.GetAtVersion(ctx, id, int(doc.CurrentVersion))
}

// GetAtVersion retrieves a document at a specific version
// 1. Get all patches up to and including the target version
// 2. Start with an empty document {}
// 3. Apply each patch in order (version 1, 2, ..., N)
// 4. Return the reconstructed document
func (s *DocumentStore) GetAtVersion(ctx context.Context, id uuid.UUID, version int) (*Document, error) {
	var docID pgtype.UUID
	setUUID(&docID, id)

	// Get document metadata
	doc, err := s.queries.GetDocument(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	if version > int(doc.CurrentVersion) || version < 1 {
		return nil, fmt.Errorf("version %d not found", version)
	}

	versions, err := s.queries.GetDocumentVersions(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document versions: %w", err)
	}

	// sanity check - versions must be in order
	for i, v := range versions {
		if int(v.Version) != i+1 {
			return nil, fmt.Errorf("invalid version history: expected version %d but got %d", i+1, v.Version)
		}
	}

	// Filter to only include versions <= requested version
	versions = versions[:version]

	// Reconstruct the document by applying patches
	content := []byte(`{}`)
	for _, v := range versions {
		var err error
		content, err = applyPatch(content, v.Patch)
		if err != nil {
			return nil, fmt.Errorf("failed to apply patch for version %d: %w", v.Version, err)
		}
	}

	var rawJson map[string]interface{}
	err = json.Unmarshal(content, &rawJson)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal content: %w", err)
	}

	return &Document{
		ID:             id,
		Name:           doc.Name,
		CurrentVersion: version,
		Content:        rawJson,
		CreatedAt:      doc.CreatedAt.Time,
	}, nil
}

// Update updates a document with new content, creating a new version
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

	forwardPatch, invertedPatch, err := createPatch(currentBytes, newBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create patch: %w", err)
	}

	newVersion := currentDoc.CurrentVersion + 1
	docID := pgtype.UUID{}
	setUUID(&docID, id)
	_, err = s.queries.CreateDocumentVersion(ctx, database.CreateDocumentVersionParams{
		DocumentID:    docID,
		Version:       int32(newVersion),
		Patch:         forwardPatch,
		InvertedPatch: invertedPatch,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document version: %w", err)
	}

	err = s.queries.UpdateDocumentVersion(ctx, database.UpdateDocumentVersionParams{
		ID:             docID,
		CurrentVersion: int32(newVersion),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update document version: %w", err)
	}

	return &Document{
		ID:             id,
		Name:           currentDoc.Name,
		CurrentVersion: newVersion,
		Content:        content,
		CreatedAt:      currentDoc.CreatedAt,
	}, nil
}

// ListVersions returns the version history for a document
// 1. Get all versions for the document
// 2. Return version metadata (version number, created_at)
func (s *DocumentStore) ListVersions(ctx context.Context, id uuid.UUID) (int, []VersionInfo, error) {
	var docID pgtype.UUID
	setUUID(&docID, id)

	// Verify document exists
	_, err := s.queries.GetDocument(ctx, docID)
	if err != nil {
		return 0, nil, fmt.Errorf("document not found: %w", err)
	}

	rows, err := s.queries.GetVersionHistory(ctx, docID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get version history: %w", err)
	}

	var versions []VersionInfo
	for _, v := range rows {
		versions = append(versions, VersionInfo{
			Version:   int(v.Version),
			CreatedAt: v.CreatedAt.Time,
		})
	}

	return len(versions), versions, nil
}

// Revert reverts a document to a specific version by creating a new version
// with the content from the target version
// 1. Get the document at the target version
// 2. Create a new version with that content (like Update)
func (s *DocumentStore) Revert(ctx context.Context, id uuid.UUID, targetVersion int) (*Document, error) {
	targetDoc, err := s.GetAtVersion(ctx, id, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get target version: %w", err)
	}

	doc, err := s.Update(ctx, id, targetDoc.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return doc, nil
}

// Helper functions for JSON patch operations

// jsondiff Patch.String() returns results like "patch1\npatch2\npatch3", but
// to store them in a jsonb column we need to convert it to a valid JSON array
// eg. ["patch1","patch2","patch3"]
func convertPatchToJsonArray(patch jsondiff.Patch) []byte {
	lines := strings.Split(patch.String(), "\n")
	return []byte(fmt.Sprintf("[%s]", strings.Join(lines, ",")))
}

// createPatch creates a JSON patch from base to target, returning both
// the forward patch and inverted patch
func createPatch(base, target []byte) ([]byte, []byte, error) {
	patch, err := jsondiff.CompareJSON(base, target, jsondiff.Invertible())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create patch: %w", err)
	}
	invertedPatch, err := patch.Invert()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to invert patch: %w", err)
	}
	return convertPatchToJsonArray(patch), convertPatchToJsonArray(invertedPatch), nil
}

// applyPatch applies a JSON patch to a document
func applyPatch(doc []byte, patchBytes []byte) ([]byte, error) {
	patch, err := jsonpatch.DecodePatch(patchBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode patch: %w", err)
	}
	modified, err := patch.Apply(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch: %w", err)
	}
	return modified, nil
}
