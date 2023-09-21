package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/infrastructure/repository"
)

//go:generate mockgen -source=document.go -destination=./mocks/document.go -package=mocks
type DocumentRepository interface {
	GetDocument(ctx context.Context, id int32) (*domain.Document, error)
	CreateDocument(ctx context.Context, doc *domain.Document) (int32, error)
	DeleteDocument(ctx context.Context, id int32) error

	GetProjectDocuments(ctx context.Context, projectID int32) ([]domain.Document, error)
	AddDocumentToProject(ctx context.Context, docID int32, projectID int32) error
	RemoveDocumentFromProject(ctx context.Context, docID int32, projectID int32) error
}

type DocumentService struct {
	filesDir    string
	repo        DocumentRepository
	projectRepo ProjectRepository
	fileRepo    FileRepository
}

func NewDocumentService(repo DocumentRepository, projectRepo ProjectRepository, fileRepo FileRepository) *DocumentService {
	mimetype.SetLimit(0) // Make mime type detector read all file
	return &DocumentService{"documents", repo, projectRepo, fileRepo}
}

func (s *DocumentService) GetDocument(ctx context.Context, id int32) (*domain.Document, error) {
	doc, err := s.repo.GetDocument(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("document_id")
		}
		return nil, fmt.Errorf("getting document %d: %w", id, err)
	}

	content, err := s.fileRepo.Read(ctx, filepath.Join(s.filesDir, doc.FileID))
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("document_id")
		}
		return nil, fmt.Errorf("reading document %d content: %w", id, err)
	}
	doc.Content = content

	return doc, nil
}

func (s *DocumentService) GetProjectDocuments(ctx context.Context, id int32) ([]domain.Document, error) {
	_, err := s.projectRepo.GetProject(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting document %d: %w", id, err)
	}

	documents, err := s.repo.GetProjectDocuments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting project %d documents: %w", id, err)
	}

	return documents, nil
}

func (s *DocumentService) AddDocumentToProject(ctx context.Context, doc *domain.Document, projectID int32) (*domain.Document, error) {
	if doc.SizeBytes > (5 << 20) { // 5MB
		return nil, apperror.NewInvalidRequest("Document size is too big.", "")
	}

	_, err := s.projectRepo.GetProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d: %w", projectID, err)
	}

	mt := mimetype.Detect(doc.Content)
	doc.MimeType = mt.String()

	doc.FileID = uuid.New().String() + mt.Extension() // Already has dot in file extension

	doc.UserID = 1 // TODO: use auth identifier in the future

	docID, err := s.repo.CreateDocument(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("creating document: %w", err)
	}

	err = s.fileRepo.Save(ctx, doc.Content, filepath.Join(s.filesDir, doc.FileID))
	if err != nil {
		return nil, fmt.Errorf("saving document: %w", err)
	}

	err = s.repo.AddDocumentToProject(ctx, docID, projectID)
	if err != nil {
		return nil, fmt.Errorf("adding document %d to project %d: %w", docID, projectID, err)
	}

	document, err := s.repo.GetDocument(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("getting document: %w", err)
	}

	return document, nil
}

func (s *DocumentService) DeleteDocumentFromProject(ctx context.Context, docID int32, projectID int32) error {
	doc, err := s.repo.GetDocument(ctx, docID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.NewNotFound("document_id")
		}
		return fmt.Errorf("getting document %d: %w", docID, err)
	}

	err = s.repo.RemoveDocumentFromProject(ctx, docID, projectID)
	if err != nil {
		return fmt.Errorf("removing document %d from project %d: %w", docID, projectID, err)
	}

	err = s.repo.DeleteDocument(ctx, docID)
	if err != nil {
		return fmt.Errorf("deleting document %d: %w", docID, err)
	}

	err = s.fileRepo.Delete(ctx, filepath.Join(s.filesDir, doc.FileID))
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.NewNotFound("document_id")
		}
		return fmt.Errorf("deleting document %d fs: %w", docID, err)
	}

	return nil
}
