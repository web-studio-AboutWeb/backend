package http

import (
	"context"
	"net/http"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/dto"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

type ProjectCategoryService interface {
	CreateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) (*domain.ProjectCategory, error)
	GetProjectCategories(ctx context.Context) ([]domain.ProjectCategory, error)
	UpdateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) error
	DeleteProjectCategory(ctx context.Context, id int16) error
}

type projectCategoryHandler struct {
	pcService ProjectCategoryService
}

func newProjectCategoryHandler(srv ProjectCategoryService) *projectCategoryHandler {
	return &projectCategoryHandler{srv}
}

// createProjectCategory godoc
// @Summary      Create project category
// @Description  Creates a new project category.
// @Tags         Project categories
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProjectCategoryRequest true "Request body."
// @Success      200  {object}  domain.ProjectCategory
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/projects/categories [post]
func (h *projectCategoryHandler) createProjectCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectCategoryRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.pcService.CreateProjectCategory(r.Context(), &domain.ProjectCategory{Name: req.Name})
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getProjectCategories godoc
// @Summary      Get project categories
// @Tags         Project categories
// @Accept       json
// @Produce      json
// @Success      200  {array} domain.ProjectCategory
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/projects/categories [get]
func (h *projectCategoryHandler) getProjectCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := h.pcService.GetProjectCategories(r.Context())
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, resp, w)
}

// updateProjectCategory godoc
// @Summary      Update project category
// @Description  Updates a project category.
// @Tags         Project categories
// @Accept       json
// @Produce      json
// @Param        category_id path int true "Project category identifier."
// @Param        request body dto.UpdateProjectCategoryRequest true "Request body."
// @Success      204
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/projects/categories/{category_id} [put]
func (h *projectCategoryHandler) updateProjectCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := httphelp.ParseParamInt16("category_id", r)

	var req dto.UpdateProjectCategoryRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	err := h.pcService.UpdateProjectCategory(r.Context(), &domain.ProjectCategory{
		ID:   categoryID,
		Name: req.Name,
	})
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// deleteProjectCategory godoc
// @Summary      Delete project category
// @Tags         Project categories
// @Param        category_id path int true "Project category identifier."
// @Success      204
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/projects/categories/{category_id} [delete]
func (h *projectCategoryHandler) deleteProjectCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := httphelp.ParseParamInt16("category_id", r)

	err := h.pcService.DeleteProjectCategory(r.Context(), categoryID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
