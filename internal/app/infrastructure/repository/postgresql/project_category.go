package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type ProjectCategoryRepository struct {
	pool Driver
}

func NewProjectCategoryRepository(pool Driver) *ProjectCategoryRepository {
	return &ProjectCategoryRepository{pool}
}

func (r *ProjectCategoryRepository) CreateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) (int16, error) {
	var id int16

	err := r.pool.QueryRow(ctx, `
		INSERT INTO project_categories(name)
		VALUES ($1)`, pc.Name).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting project category: %w", err)
	}

	return id, nil
}

func (r *ProjectCategoryRepository) GetProjectCategories(ctx context.Context) ([]domain.ProjectCategory, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name
		FROM project_categories
		ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("selecting project categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.ProjectCategory
	var pc domain.ProjectCategory

	for rows.Next() {
		err = rows.Scan(&pc.ID, &pc.Name)
		if err != nil {
			return nil, fmt.Errorf("scanning project category: %w", err)
		}

		categories = append(categories, pc)
	}

	return categories, nil
}

func (r *ProjectCategoryRepository) GetProjectCategory(ctx context.Context, id int16) (*domain.ProjectCategory, error) {
	var pc domain.ProjectCategory

	err := r.pool.QueryRow(ctx, `
		SELECT id, name
		FROM project_categories
		WHERE id=$1`, id).Scan(
		&pc.ID,
		&pc.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("selecting project category: %w", err)
	}

	return &pc, nil
}

func (r *ProjectCategoryRepository) GetProjectCategoryByName(ctx context.Context, name string) (*domain.ProjectCategory, error) {
	var pc domain.ProjectCategory

	err := r.pool.QueryRow(ctx, `
		SELECT id, name
		FROM project_categories
		WHERE lower(name)=lower($1)`, name).Scan(
		&pc.ID,
		&pc.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("selecting project category: %w", err)
	}

	return &pc, nil
}

func (r *ProjectCategoryRepository) UpdateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE project_categories
		SET name=$2
		WHERE id=$1`, pc.ID, pc.Name)
	if err != nil {
		return fmt.Errorf("updating project category: %w", err)
	}

	return nil
}

func (r *ProjectCategoryRepository) DeleteProjectCategory(ctx context.Context, id int16) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM project_categories
		WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting project category: %w", err)
	}

	return nil
}
