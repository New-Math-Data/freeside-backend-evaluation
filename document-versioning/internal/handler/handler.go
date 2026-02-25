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
// TODO: Implement this handler
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

	// TODO: Call the store to create the document
	// doc, err := h.store.Create(c.Request.Context(), req.Name, req.Content)

	c.JSON(http.StatusNotImplemented, api.ErrorResponse{
		Error: "not implemented - complete this handler",
	})
}

// GetDocument retrieves a document, optionally at a specific version
// TODO: Implement this handler
// 1. The document ID is already parsed by the generated code
// 2. Check if a version query parameter was provided (params.Version)
// 3. If version provided, call h.store.GetAtVersion()
// 4. If no version, call h.store.GetCurrent()
// 5. Return the document as api.DocumentResponse
// 6. Handle errors (404 for not found)
func (h *Handler) GetDocument(c *gin.Context, id openapi_types.UUID, params api.GetDocumentParams) {
	// The ID is already parsed by the generated wrapper code
	// Use id directly (it's a uuid.UUID under the hood)

	_ = id // Use this parsed ID

	// TODO: Implement document retrieval
	// Check params.Version to see if a specific version was requested

	c.JSON(http.StatusNotImplemented, api.ErrorResponse{
		Error: "not implemented - complete this handler",
	})
}

// UpdateDocument updates a document, creating a new version
// TODO: Implement this handler
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

	_ = id // Use this parsed ID

	// TODO: Call the store to update the document

	c.JSON(http.StatusNotImplemented, api.ErrorResponse{
		Error: "not implemented - complete this handler",
	})
}

// ListVersions returns the version history for a document
// TODO: Implement this handler
// 1. The document ID is already parsed by the generated code
// 2. Call h.store.ListVersions()
// 3. Return the version list as api.VersionListResponse
// 4. Handle errors (404 for not found)
func (h *Handler) ListVersions(c *gin.Context, id openapi_types.UUID) {
	_ = id // Use this parsed ID

	// TODO: Call the store to list versions

	c.JSON(http.StatusNotImplemented, api.ErrorResponse{
		Error: "not implemented - complete this handler",
	})
}

// RevertDocument reverts a document to a specific version
// TODO: Implement this handler
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

	_ = id // Use this parsed ID

	// TODO: Call the store to revert the document

	c.JSON(http.StatusNotImplemented, api.ErrorResponse{
		Error: "not implemented - complete this handler",
	})
}
