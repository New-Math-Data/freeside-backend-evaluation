package store

import (
	"context"
	"errors"

	"document-versioning/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
)

// MockDatabaseQueries is a mock implementation of database.Queries for testing
type MockDatabaseQueries struct {
	CreateDocumentFunc        func(ctx context.Context, name string) (database.Document, error)
	CreateDocumentVersionFunc func(ctx context.Context, arg database.CreateDocumentVersionParams) (database.DocumentVersion, error)
	GetDocumentFunc           func(ctx context.Context, id pgtype.UUID) (database.Document, error)
	GetDocumentVersionsFunc   func(ctx context.Context, documentID pgtype.UUID) ([]database.DocumentVersion, error)
	UpdateDocumentVersionFunc func(ctx context.Context, arg database.UpdateDocumentVersionParams) error
	GetVersionHistoryFunc     func(ctx context.Context, documentID pgtype.UUID) ([]database.GetVersionHistoryRow, error)
}

func (m *MockDatabaseQueries) CreateDocument(ctx context.Context, name string) (database.Document, error) {
	if m.CreateDocumentFunc != nil {
		return m.CreateDocumentFunc(ctx, name)
	}
	return database.Document{}, errors.New("not implemented")
}

func (m *MockDatabaseQueries) CreateDocumentVersion(ctx context.Context, arg database.CreateDocumentVersionParams) (database.DocumentVersion, error) {
	if m.CreateDocumentVersionFunc != nil {
		return m.CreateDocumentVersionFunc(ctx, arg)
	}
	return database.DocumentVersion{}, errors.New("not implemented")
}

func (m *MockDatabaseQueries) GetDocument(ctx context.Context, id pgtype.UUID) (database.Document, error) {
	if m.GetDocumentFunc != nil {
		return m.GetDocumentFunc(ctx, id)
	}
	return database.Document{}, errors.New("not implemented")
}

func (m *MockDatabaseQueries) GetDocumentVersions(ctx context.Context, documentID pgtype.UUID) ([]database.DocumentVersion, error) {
	if m.GetDocumentVersionsFunc != nil {
		return m.GetDocumentVersionsFunc(ctx, documentID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDatabaseQueries) UpdateDocumentVersion(ctx context.Context, arg database.UpdateDocumentVersionParams) error {
	if m.UpdateDocumentVersionFunc != nil {
		return m.UpdateDocumentVersionFunc(ctx, arg)
	}
	return errors.New("not implemented")
}

func (m *MockDatabaseQueries) GetVersionHistory(ctx context.Context, documentID pgtype.UUID) ([]database.GetVersionHistoryRow, error) {
	if m.GetVersionHistoryFunc != nil {
		return m.GetVersionHistoryFunc(ctx, documentID)
	}
	return nil, errors.New("not implemented")
}
