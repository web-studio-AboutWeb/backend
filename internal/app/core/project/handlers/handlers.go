package handlers

import (
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type ProjectHandlers struct {
	CreateProjectHandler      *CreateProjectHandler
	GetProjectHandler         *GetProjectHandler
	GetProjectsHandler        *GetProjectsHandler
	GetProjectStaffersHandler *GetProjectParticipantsHandler
	UpdateProjectHandler      *UpdateProjectHandler
	DeleteProjectHandler      *DeleteProjectHandler
}

func New(gateways *gateways.Gateways) (*ProjectHandlers, error) {
	getProjectHandler := NewGetProjectHandler(gateways.ProjectGateway)
	getProjectsHandler := NewGetProjectsHandler(gateways.ProjectGateway)
	getProjectParticipantsHandler := NewGetProjectParticipantsHandler(
		gateways.ProjectGateway, getProjectHandler,
	)
	createProjectHandler := NewCreateProjectHandler(
		gateways.ProjectGateway, getProjectHandler,
	)
	updateProjectHandler := NewUpdateProjectHandler(
		gateways.ProjectGateway, getProjectHandler,
	)
	deleteProjectHandler := NewDeleteProjectHandler(
		gateways.ProjectGateway, getProjectHandler,
	)

	return &ProjectHandlers{
		CreateProjectHandler:      createProjectHandler,
		GetProjectHandler:         getProjectHandler,
		GetProjectsHandler:        getProjectsHandler,
		GetProjectStaffersHandler: getProjectParticipantsHandler,
		UpdateProjectHandler:      updateProjectHandler,
		DeleteProjectHandler:      deleteProjectHandler,
	}, nil
}
