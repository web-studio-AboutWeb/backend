package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

//go:generate mockgen -source=document.go -destination=./mocks/document.go -package=mocks
type DocumentService interface {
	GetDocument(ctx context.Context, id int32) (*domain.Document, error)

	GetProjectDocuments(ctx context.Context, id int32) ([]domain.Document, error)
	AddDocumentToProject(ctx context.Context, doc *domain.Document, projectID int32) (*domain.Document, error)
	DeleteDocumentFromProject(ctx context.Context, docID int32, projectID int32) error
}

type documentHandler struct {
	documentService DocumentService
}

func newDocumentHandler(documentService DocumentService) *documentHandler {
	return &documentHandler{documentService}
}

// downloadDocument godoc
// @Summary      Download document
// @Description  Returns document file content if document exists.
// @Tags         Documents
// @Produce      octet-stream
// @Param        document_id path int true "Document identifier."
// @Success      200
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/documents/{document_id} [get]
func (h *documentHandler) downloadDocument(w http.ResponseWriter, r *http.Request) {
	did := httphelp.ParseParamInt32("document_id", r)

	doc, err := h.documentService.GetDocument(r.Context(), did)
	if err != nil {
		httphelp.SendError(fmt.Errorf("getting document: %w", err), w)
		return
	}

	http.ServeContent(w, r, doc.OriginalFilename, doc.CreatedAt, bytes.NewReader(doc.Content))
}

// getProjectDocuments godoc
// @Summary      Get project documents
// @Description  Returns list of project documents.
// @Tags         Documents
// @Param        project_id path int true "Project identifier."
// @Success      200  {array}   domain.Document
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/documents [get]
func (h *documentHandler) getProjectDocuments(w http.ResponseWriter, r *http.Request) {
	pid := httphelp.ParseParamInt32("project_id", r)

	docs, err := h.documentService.GetProjectDocuments(r.Context(), pid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("getting project documents: %w", err), w)
		return
	}

	httphelp.SendJSON(http.StatusOK, docs, w)
}

// addDocumentToProject godoc
// @Security     CSRF
// @Summary      Add document to project
// @Description  Adds document to project.
// @Description
// @Description  Accepts `multipart/form-data` and document file.
// @Tags         Documents
// @Accept       mpfd
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Param        file formData file true "Document file."
// @Success      200  {object}	domain.Document
// @Failure      400  {object}  apperror.Error
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/documents [post]
func (h *documentHandler) addDocumentToProject(w http.ResponseWriter, r *http.Request) {
	pid := httphelp.ParseParamInt32("project_id", r)

	file, header, err := r.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			// TODO: custom http error
			httphelp.SendError(fmt.Errorf("file is not presented"), w)
			return
		}
		httphelp.SendError(fmt.Errorf("parsing form file: %w", err), w)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		httphelp.SendError(fmt.Errorf("reading file: %w", err), w)
		return
	}

	doc := &domain.Document{
		OriginalFilename: header.Filename,
		Size:             int32(header.Size),
		Content:          content,
	}

	document, err := h.documentService.AddDocumentToProject(r.Context(), doc, pid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("adding document to project: %w", err), w)
		return
	}

	httphelp.SendJSON(http.StatusOK, document, w)
}

// removeDocumentFromProject godoc
// @Summary      Delete document from project
// @Description  Deletes document from projects.
// @Description
// @Description  **Also deletes file forever, it will never be accepted in the future.**
// @Tags         Documents
// @Param        project_id path int true "Project identifier."
// @Param        document_id path int true "Document identifier."
// @Success      200
// @Failure      400  {object}  apperror.Error
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/documents/{document_id} [delete]
func (h *documentHandler) removeDocumentFromProject(w http.ResponseWriter, r *http.Request) {
	pid := httphelp.ParseParamInt32("project_id", r)
	did := httphelp.ParseParamInt32("document_id", r)

	err := h.documentService.DeleteDocumentFromProject(r.Context(), did, pid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("deleting document from project: %w", err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
