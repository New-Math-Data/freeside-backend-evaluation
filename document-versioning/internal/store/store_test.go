package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"document-versioning/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions

func createPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func createPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

// Tests for Create method

func TestCreate_Successful(t *testing.T) {
	testID := uuid.New()
	createdAt := time.Now()
	testContent := map[string]interface{}{
		"title": "Test Document",
		"text":  "Some content",
	}

	mock := &MockDatabaseQueries{
		CreateDocumentFunc: func(ctx context.Context, name string) (database.Document, error) {
			assert.Equal(t, "My Document", name)
			return database.Document{
				ID:             createPgUUID(testID),
				Name:           name,
				CurrentVersion: 1,
				CreatedAt:      createPgTimestamp(createdAt),
			}, nil
		},
		CreateDocumentVersionFunc: func(ctx context.Context, arg database.CreateDocumentVersionParams) (database.DocumentVersion, error) {
			// Verify the version was created with correct parameters
			assert.Equal(t, int32(1), arg.Version)
			assert.NotEmpty(t, arg.Patch, "patch should not be empty")
			assert.NotEmpty(t, arg.InvertedPatch, "inverted patch should not be empty")

			// Verify that the patch is valid JSON
			var patches []interface{}
			err := json.Unmarshal(arg.Patch, &patches)
			assert.NoError(t, err, "patch should be valid JSON")
			assert.NotEmpty(t, patches, "patch should contain operations")

			// Verify inverted patch is also valid JSON
			var invertedPatches []interface{}
			err = json.Unmarshal(arg.InvertedPatch, &invertedPatches)
			assert.NoError(t, err, "inverted patch should be valid JSON")

			return database.DocumentVersion{
				DocumentID:    arg.DocumentID,
				Version:       arg.Version,
				Patch:         arg.Patch,
				InvertedPatch: arg.InvertedPatch,
				CreatedAt:     createPgTimestamp(createdAt),
			}, nil
		},
	}

	store := NewDocumentStore(mock)
	doc, err := store.Create(context.Background(), "My Document", testContent)

	require.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, testID, doc.ID)
	assert.Equal(t, "My Document", doc.Name)
	assert.Equal(t, 1, doc.CurrentVersion)
	assert.Equal(t, testContent, doc.Content)
	assert.Equal(t, createdAt, doc.CreatedAt)
}

func TestCreate_CreateDocument_FailsOnDatabaseError(t *testing.T) {
	mock := &MockDatabaseQueries{
		CreateDocumentFunc: func(ctx context.Context, name string) (database.Document, error) {
			return database.Document{}, errors.New("database connection failed")
		},
	}

	store := NewDocumentStore(mock)
	doc, err := store.Create(context.Background(), "My Document", map[string]interface{}{})

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Contains(t, err.Error(), "failed to create document")
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestCreate_CreateDocumentVersion_FailsOnDatabaseError(t *testing.T) {
	testID := uuid.New()
	createdAt := time.Now()
	testContent := map[string]interface{}{"title": "Test"}

	mock := &MockDatabaseQueries{
		CreateDocumentFunc: func(ctx context.Context, name string) (database.Document, error) {
			return database.Document{
				ID:             createPgUUID(testID),
				Name:           name,
				CurrentVersion: 1,
				CreatedAt:      createPgTimestamp(createdAt),
			}, nil
		},
		CreateDocumentVersionFunc: func(ctx context.Context, arg database.CreateDocumentVersionParams) (database.DocumentVersion, error) {
			return database.DocumentVersion{}, errors.New("version storage failed")
		},
	}

	store := NewDocumentStore(mock)
	doc, err := store.Create(context.Background(), "My Document", testContent)

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Contains(t, err.Error(), "failed to create document version")
	assert.Contains(t, err.Error(), "version storage failed")
}

func TestCreate_WithEmptyContent(t *testing.T) {
	testID := uuid.New()
	createdAt := time.Now()
	testContent := map[string]interface{}{}

	mock := &MockDatabaseQueries{
		CreateDocumentFunc: func(ctx context.Context, name string) (database.Document, error) {
			return database.Document{
				ID:             createPgUUID(testID),
				Name:           name,
				CurrentVersion: 1,
				CreatedAt:      createPgTimestamp(createdAt),
			}, nil
		},
		CreateDocumentVersionFunc: func(ctx context.Context, arg database.CreateDocumentVersionParams) (database.DocumentVersion, error) {
			return database.DocumentVersion{
				DocumentID:    arg.DocumentID,
				Version:       arg.Version,
				Patch:         arg.Patch,
				InvertedPatch: arg.InvertedPatch,
				CreatedAt:     createPgTimestamp(createdAt),
			}, nil
		},
	}

	store := NewDocumentStore(mock)
	doc, err := store.Create(context.Background(), "Empty Doc", testContent)

	require.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, testID, doc.ID)
	assert.Equal(t, testContent, doc.Content)
	assert.Equal(t, 1, doc.CurrentVersion)
}

func TestCreate_WithSampleContent(t *testing.T) {
	testID := uuid.New()
	createdAt := time.Now()

	// Load sample document from JSON file
	exe, err := os.Executable()
	require.NoError(t, err, "failed to read sample_document.json")
	exeDir := path.Dir(exe)
	sampleData, err := os.ReadFile(path.Join(exeDir, "../../../testdata/sample_document.json"))
	require.NoError(t, err, "failed to read sample_document.json")

	var testContent map[string]interface{}
	err = json.Unmarshal(sampleData, &testContent)
	require.NoError(t, err, "failed to unmarshal sample_document.json")

	mock := &MockDatabaseQueries{
		CreateDocumentFunc: func(ctx context.Context, name string) (database.Document, error) {
			return database.Document{
				ID:             createPgUUID(testID),
				Name:           name,
				CurrentVersion: 1,
				CreatedAt:      createPgTimestamp(createdAt),
			}, nil
		},
		CreateDocumentVersionFunc: func(ctx context.Context, arg database.CreateDocumentVersionParams) (database.DocumentVersion, error) {
			// Verify the version was created with correct parameters
			assert.Equal(t, int32(1), arg.Version)
			assert.NotEmpty(t, arg.Patch, "patch should not be empty")
			assert.NotEmpty(t, arg.InvertedPatch, "inverted patch should not be empty")

			// Verify that the patch is valid JSON
			var patches []interface{}
			err := json.Unmarshal(arg.Patch, &patches)
			assert.NoError(t, err, "patch should be valid JSON")
			assert.NotEmpty(t, patches, "patch should contain operations")

			// Verify inverted patch is also valid JSON
			var invertedPatches []interface{}
			err = json.Unmarshal(arg.InvertedPatch, &invertedPatches)
			assert.NoError(t, err, "inverted patch should be valid JSON")

			return database.DocumentVersion{
				DocumentID:    arg.DocumentID,
				Version:       arg.Version,
				Patch:         arg.Patch,
				InvertedPatch: arg.InvertedPatch,
				CreatedAt:     createPgTimestamp(createdAt),
			}, nil
		},
	}

	store := NewDocumentStore(mock)
	doc, err := store.Create(context.Background(), "Infrastructure Asset", testContent)

	require.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, testContent, doc.Content)

	// Verify the nested structure is preserved
	assert.Equal(t, "Infrastructure Asset", doc.Content["title"])
	assert.Equal(t, "electrical_transformer", doc.Content["type"])
	assert.Equal(t, "operational", doc.Content["status"])

	// Verify location structure
	location := doc.Content["location"].(map[string]interface{})
	assert.Equal(t, 40.7128, location["latitude"])
	assert.Equal(t, -74.0060, location["longitude"])
	assert.Equal(t, "123 Main Street", location["address"])

	// Verify tags array
	tags := doc.Content["tags"].([]interface{})
	assert.Len(t, tags, 3)
	assert.Equal(t, "critical", tags[0])
}
