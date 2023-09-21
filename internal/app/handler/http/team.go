package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/dto"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

type TeamService interface {
	GetTeam(ctx context.Context, id int32) (*domain.Team, error)
	GetTeams(ctx context.Context) ([]domain.Team, error)
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	UpdateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	SetTeamImage(ctx context.Context, teamID int32, img []byte) error
	GetTeamImage(ctx context.Context, teamID int32) (*domain.Team, error)
	DisableTeam(ctx context.Context, teamID int32) error
	EnableTeam(ctx context.Context, teamID int32) error
}

type teamHandler struct {
	teamService TeamService
}

func newTeamHandler(service TeamService) *teamHandler {
	return &teamHandler{service}
}

// getTeam godoc
// @Summary      Get team by identifier
// @Description  Returns information about single team.
// @Tags         Teams
// @Produce      json
// @Param        team_id path int true "Team identifier."
// @Success      200  {object}  domain.Team
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams/{team_id} [get]
func (h *teamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("team_id", r)

	response, err := h.teamService.GetTeam(r.Context(), tid)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getTeams godoc
// @Summary      Get teams
// @Description  Returns list of teams.
// @Tags         Teams
// @Produce      json
// @Success      200  {array}  domain.Team
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams [get]
func (h *teamHandler) getTeams(w http.ResponseWriter, r *http.Request) {
	response, err := h.teamService.GetTeams(r.Context())
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// createTeam godoc
// @Summary      Create team
// @Description  Creates a new team. Returns an object with information about created team.
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTeamRequest true "Request body."
// @Success      200  {object}	domain.Team
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams [post]
func (h *teamHandler) createTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTeamRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.teamService.CreateTeam(r.Context(), req.ToDomain())
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// updateTeam godoc
// @Summary      Update team
// @Description  Updates a team.
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        team_id path int true "Team identifier."
// @Param        request body dto.UpdateTeamRequest true "Request body."
// @Success      200  {object}	domain.Team
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams/{team_id} [put]
func (h *teamHandler) updateTeam(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("team_id", r)

	var req dto.UpdateTeamRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.teamService.UpdateTeam(r.Context(), req.ToDomain(tid))
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// setTeamImage godoc
// @Summary      Set team image
// @Description  Updated team image. Accepts `multipart/form-data`.
// @Description
// @Description  Note: if a team already has an image, it will be deleted automatically on success.
// @Tags         Teams
// @Accept       mpfd
// @Param        team_id path int true "Team identifier."
// @Param        file formData file true "Image file. MUST have one of the following mime types: [`image/jpeg`, `image/png`, `image/webp`]"
// @Success      200
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams/{team_id}/image [post]
func (h *teamHandler) setTeamImage(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("team_id", r)

	file, _, err := r.FormFile("file")
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

	err = h.teamService.SetTeamImage(r.Context(), tid, content)
	if err != nil {
		httphelp.SendError(fmt.Errorf("setting team image: %w", err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getTeamImage godoc
// @Summary      Get team image content
// @Description  Returns team image.
// @Tags         Teams
// @Produce      octet-stream
// @Param        team_id path int true "Team identifier."
// @Success      200
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams/{team_id}/image [get]
func (h *teamHandler) getTeamImage(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("team_id", r)

	response, err := h.teamService.GetTeamImage(r.Context(), tid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("getting team image: %w", err), w)
		return
	}

	fileName := fmt.Sprintf("%s.%s", response.Title, filepath.Ext(response.ImageID))

	http.ServeContent(w, r, fileName, response.UpdatedAt, bytes.NewReader(response.ImageContent))
}

// disableTeam godoc
// @Summary      Disable team
// @Description  Disables a team.
// @Tags         Teams
// @Param        team_id path int true "Team identifier."
// @Success      200
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams/{team_id}/disable [post]
func (h *teamHandler) disableTeam(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("team_id", r)

	err := h.teamService.DisableTeam(r.Context(), tid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("disabling team: %w", err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// enableTeam godoc
// @Summary      Enable team
// @Description  Enables a team.
// @Tags         Teams
// @Param        team_id path int true "Team identifier."
// @Success      200
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/teams/{team_id}/enable [post]
func (h *teamHandler) enableTeam(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("team_id", r)

	err := h.teamService.EnableTeam(r.Context(), tid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("enabling team: %w", err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
