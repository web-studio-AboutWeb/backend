package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type DocumentRepository struct {
	pool Driver
}

func NewDocumentRepository(pool Driver) *DocumentRepository {
	return &DocumentRepository{pool}
}

func (r *DocumentRepository) GetDocument(ctx context.Context, id int32) (*domain.Document, error) {
	var doc domain.Document

	err := r.pool.QueryRow(ctx, `
		SELECT id, filename, file_id, mime, size, user_id, created_at
		FROM documents
		WHERE id=$1`, id).Scan(
		&doc.ID,
		&doc.OriginalFilename,
		&doc.FileID,
		&doc.MimeType,
		&doc.Size,
		&doc.UserID,
		&doc.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning document: %w", err)
	}

	return &doc, nil
}

func (r *DocumentRepository) CreateDocument(ctx context.Context, doc *domain.Document) (int32, error) {
	var id int32

	err := r.pool.QueryRow(ctx, `
		INSERT INTO documents(filename, file_id, mime, size, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		doc.OriginalFilename,
		doc.FileID,
		doc.MimeType,
		doc.Size,
		doc.UserID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting document: %w", err)
	}

	return id, nil
}

func (r *DocumentRepository) DeleteDocument(ctx context.Context, id int32) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM documents WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting document: %w", err)
	}

	return nil
}

func (r *DocumentRepository) GetProjectDocuments(ctx context.Context, projectID int32) ([]domain.Document, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT d.id, d.filename, d.file_id, d.mime, d.size, d.user_id, d.created_at
		FROM documents d
		JOIN project_documents pd ON pd.document_id=d.id
		WHERE pd.project_id=$1`, projectID)
	if err != nil {
		return nil, fmt.Errorf("getting project documents: %w", err)
	}
	defer rows.Close()

	var docs []domain.Document
	for rows.Next() {
		var doc domain.Document

		err = rows.Scan(
			&doc.ID,
			&doc.OriginalFilename,
			&doc.FileID,
			&doc.MimeType,
			&doc.Size,
			&doc.UserID,
			&doc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning document: %w", err)
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

func (r *DocumentRepository) AddDocumentToProject(ctx context.Context, docID, projectID int32) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO project_documents(document_id, project_id)
		VALUES ($1, $2)`, docID, projectID)
	if err != nil {
		return fmt.Errorf("inserting project document: %w", err)
	}

	return nil
}

func (r *DocumentRepository) RemoveDocumentFromProject(ctx context.Context, docID, projectID int32) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM project_documents
		WHERE document_id=$1 AND project_id=$2`, docID, projectID)
	if err != nil {
		return fmt.Errorf("deleting project document: %w", err)
	}

	return nil
}
