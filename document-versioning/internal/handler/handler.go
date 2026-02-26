package handler

import (
	"net/http"

	"document-versioning/internal/api"
	"document-versioning/internal/store"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Handler implements the api.ServerInterface
type Handler struct {
	store *store.DocumentStore
}

// NewHandler creates a new Handler
func NewHandler(store *store.DocumentStore) *Handler {
	return &Handler{
		store: store,
	}
}

// Ensure Handler implements ServerInterface
var _ api.ServerInterface = (*Handler)(nil)

// HealthCheck returns the health status of the service
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, api.HealthResponse{
		Status: "ok",
	})
}

// CreateDocument creates a new document
// 1. Parse the request body into api.CreateDocumentRequest
// 2. Call h.store.Create() with the name and content
// 3. Return the created document as api.DocumentResponse with status 201
// 4. Handle errors appropriately (400 for bad request, 500 for server error)
func (h *Handler) CreateDocument(c *gin.Context) {
	var req api.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: "name is required",
		})
		return
	}
	if req.Content == nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: "content is required",
		})
		return
	}

	doc, err := h.store.Create(c.Request.Context(), req.Name, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: "failed to create document: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, api.DocumentResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Content:   doc.Content,
		CreatedAt: doc.CreatedAt,
	})
}

// GetDocument retrieves a document, optionally at a specific version
// 1. The document ID is already parsed by the generated code
// 2. Check if a version query parameter was provided (params.Version)
// 3. If version provided, call h.store.GetAtVersion()
// 4. If no version, call h.store.GetCurrent()
// 5. Return the document as api.DocumentResponse
// 6. Handle errors (404 for not found)
func (h *Handler) GetDocument(c *gin.Context, id openapi_types.UUID, params api.GetDocumentParams) {
	var doc *store.Document
	var err error
	var version int

	if params.Version != nil {
		doc, err = h.store.GetAtVersion(c.Request.Context(), id, *params.Version)
		version = *params.Version
	} else {
		doc, err = h.store.GetCurrent(c.Request.Context(), id)
		version = doc.CurrentVersion
	}
	if err != nil {
		c.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: "document not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.DocumentResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Content:   doc.Content,
		CreatedAt: doc.CreatedAt,
		Version:   version,
	})
}

// UpdateDocument updates a document, creating a new version
// 1. The document ID is already parsed by the generated code
// 2. Parse the request body into api.UpdateDocumentRequest
// 3. Call h.store.Update() with the ID and new content
// 4. Return the updated document as api.DocumentResponse
// 5. Handle errors (404 for not found, 400 for bad request)
func (h *Handler) UpdateDocument(c *gin.Context, id openapi_types.UUID) {
	var req api.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	doc, err := h.store.Update(c.Request.Context(), id, req.Content)
	if err != nil {
		c.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: "document not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.DocumentResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Content:   doc.Content,
		CreatedAt: doc.CreatedAt,
		Version:   doc.CurrentVersion,
	})
}

// ListVersions returns the version history for a document
// 1. The document ID is already parsed by the generated code
// 2. Call h.store.ListVersions()
// 3. Return the version list as api.VersionListResponse
// 4. Handle errors (404 for not found)
func (h *Handler) ListVersions(c *gin.Context, id openapi_types.UUID) {
	count, storedVersions, err := h.store.ListVersions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: "document not found: " + err.Error(),
		})
		return
	}

	var versions []api.VersionInfo
	for _, v := range storedVersions {
		versions = append(versions, api.VersionInfo{
			Version:   v.Version,
			CreatedAt: v.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, api.VersionListResponse{
		CurrentVersion: count,
		DocumentId:     id,
		Versions:       versions,
	})
}

// RevertDocument reverts a document to a specific version
// 1. The document ID is already parsed by the generated code
// 2. Parse the request body into api.RevertRequest
// 3. Call h.store.Revert() with the ID and target version
// 4. Return the reverted document as api.DocumentResponse
// 5. Handle errors (404 for not found, 400 for invalid version)
func (h *Handler) RevertDocument(c *gin.Context, id openapi_types.UUID) {
	var req api.RevertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	doc, err := h.store.Revert(c.Request.Context(), id, req.Version)
	if err != nil {
		c.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: "document not found or version invalid: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.DocumentResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Content:   doc.Content,
		CreatedAt: doc.CreatedAt,
		Version:   doc.CurrentVersion,
	})
}
