package postgresql_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"
	"web-studio-backend/internal/app/infrastructure/repository/postgresql"
	"web-studio-backend/internal/pkg/ptr"
)

func prepareProjectMock(t *testing.T) (pgxmock.PgxPoolIface, *postgresql.ProjectRepository) {
	t.Helper()

	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	repo := postgresql.NewProjectRepository(mock)

	return mock, repo
}

func TestProjectRepository_GetProject(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		SELECT 
		    id, title, description, cover_id, created_at, updated_at, ended_at,
		    link, isactive, technologies, team_id
        FROM projects
        WHERE id = $1`

	tempTime := time.Now()

	tests := []struct {
		name        string
		response    *domain.Project
		in          int32
		wantErr     bool
		expectedErr error
		mock        func(id int32)
	}{
		{
			name: "should pass",
			in:   1,
			response: &domain.Project{
				ID:           1,
				Title:        "title",
				Description:  "description",
				IsActive:     true,
				Technologies: []string{"tech1", "tech2"},
				CreatedAt:    tempTime,
				UpdatedAt:    tempTime,
				CoverId:      ptr.String("cover_id"),
				Link:         ptr.String("link"),
				TeamID:       ptr.Int32(2),
				EndedAt:      &tempTime,
			},
			wantErr: false,
			mock: func(id int32) {
				row := mock.NewRows([]string{
					"id", "title", "description", "cover_id", "created_at", "updated_at", "ended_at",
					"link", "isactive", "technologies", "team_id",
				}).
					AddRow(
						id,
						"title",
						"description",
						ptr.String("cover_id"),
						tempTime,
						tempTime,
						&tempTime,
						ptr.String("link"),
						true,
						[]string{"tech1", "tech2"},
						ptr.Int32(2),
					)

				mock.ExpectQuery(q).WithArgs(id).WillReturnRows(row)
			},
		},
		{
			name:    "query error",
			in:      1,
			wantErr: true,
			mock: func(id int32) {
				mock.ExpectQuery(q).WithArgs(id).WillReturnError(errors.New("some error"))
			},
		},
		{
			name:        "no rows error",
			in:          1,
			wantErr:     true,
			expectedErr: repository.ErrObjectNotFound,
			mock: func(id int32) {
				mock.ExpectQuery(q).WithArgs(id).WillReturnError(pgx.ErrNoRows)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in)

			resp, err := repo.GetProject(context.Background(), tc.in)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)

				if tc.expectedErr != nil {
					require.ErrorIs(tt, err, tc.expectedErr)
				}
				return
			}

			require.NoError(tt, err)
			require.Equal(tt, tc.response, resp)
		})
	}
}

func TestProjectRepository_GetProjects(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		SELECT 
		    id, title, description, cover_id, created_at, updated_at, ended_at,
		    link, isactive, technologies, team_id
        FROM projects
        WHERE isactive
        ORDER BY created_at`

	tempTime := time.Now()

	tests := []struct {
		name     string
		response []domain.Project
		wantErr  bool
		mock     func()
	}{
		{
			name: "should pass",
			response: []domain.Project{
				{
					ID:           1,
					Title:        "title1",
					Description:  "description1",
					IsActive:     true,
					Technologies: []string{"tech1", "tech2"},
					CreatedAt:    tempTime,
					UpdatedAt:    tempTime,
					CoverId:      ptr.String("cover1"),
					Link:         ptr.String("link1"),
					TeamID:       ptr.Int32(2),
					EndedAt:      &tempTime,
				},
				{
					ID:          2,
					Title:       "title2",
					Description: "description2",
					IsActive:    false,
					CreatedAt:   tempTime,
					UpdatedAt:   tempTime,
				},
			},
			wantErr: false,
			mock: func() {
				rows := mock.NewRows([]string{
					"id", "title", "description", "cover_id", "created_at", "updated_at", "ended_at",
					"link", "isactive", "technologies", "team_id",
				}).
					AddRow(
						int32(1),
						"title1",
						"description1",
						ptr.String("cover1"),
						tempTime,
						tempTime,
						&tempTime,
						ptr.String("link1"),
						true,
						[]string{"tech1", "tech2"},
						ptr.Int32(2),
					).
					AddRow(
						int32(2),
						"title2",
						"description2",
						nil,
						tempTime,
						tempTime,
						nil,
						nil,
						false,
						nil,
						nil,
					)

				mock.ExpectQuery(q).WillReturnRows(rows).RowsWillBeClosed()
			},
		},
		{
			name:     "scan error",
			response: nil,
			wantErr:  true,
			mock: func() {
				mock.ExpectQuery(q).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock()

			resp, err := repo.GetProjects(context.Background())

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
			require.Equal(tt, tc.response, resp)
		})
	}
}

func TestProjectRepository_CreateProject(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		INSERT INTO projects(title, description, team_id, isactive, link, technologies)
		VALUES($1, $2, $3, TRUE, $4, $5)
		RETURNING id`

	tests := []struct {
		name     string
		response int32
		wantErr  bool
		mock     func(project *domain.Project, id int32)
	}{
		{
			name:     "should pass",
			response: 1,
			wantErr:  false,
			mock: func(project *domain.Project, id int32) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(id)

				mock.
					ExpectQuery(q).
					WithArgs(
						project.Title,
						project.Description,
						project.TeamID,
						project.Link,
						project.Technologies,
					).WillReturnRows(rows)
			},
		},
		{
			name:     "query error",
			response: 0,
			wantErr:  true,
			mock: func(project *domain.Project, id int32) {
				mock.
					ExpectQuery(q).
					WithArgs(
						project.Title,
						project.Description,
						project.TeamID,
						project.Link,
						project.Technologies,
					).WillReturnError(fmt.Errorf("some error"))
			},
		},
	}

	testProject := &domain.Project{
		Title:        "title",
		Description:  "description",
		TeamID:       ptr.Int32(1),
		Link:         ptr.String("test"),
		Technologies: []string{"tech1", "tech2"},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(testProject, tc.response)

			resp, err := repo.CreateProject(context.Background(), testProject)

			require.NoError(tt, mock.ExpectationsWereMet())

			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
			require.Equal(tt, tc.response, resp)

		})
	}
}

func TestProjectRepository_UpdateProject(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		UPDATE projects
		SET title=$2, description=$3, link=$4, technologies=$5, updated_at=now()
		WHERE id = $1`

	tests := []struct {
		name    string
		in      *domain.Project
		wantErr bool
		mock    func(project *domain.Project)
	}{
		{
			name:    "should pass",
			wantErr: false,
			in: &domain.Project{
				ID:           1,
				Title:        "title",
				Description:  "description",
				Link:         ptr.String("link"),
				Technologies: []string{"tech1", "tech2"},
			},
			mock: func(project *domain.Project) {
				mock.ExpectExec(q).WithArgs(
					project.ID,
					project.Title,
					project.Description,
					project.Link,
					project.Technologies,
				).WillReturnResult(pgxmock.NewResult("UPDATED", 1))
			},
		},
		{
			name:    "exec error",
			wantErr: true,
			in:      &domain.Project{},
			mock: func(project *domain.Project) {
				mock.ExpectExec(q).WithArgs(
					project.ID,
					project.Title,
					project.Description,
					project.Link,
					project.Technologies,
				).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in)

			err := repo.UpdateProject(context.Background(), tc.in)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestProjectRepository_DisableProject(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `UPDATE projects SET isactive=FALSE WHERE id = $1`

	tests := []struct {
		name    string
		id      int32
		wantErr bool
		mock    func(id int32)
	}{
		{
			name:    "should pass",
			wantErr: false,
			mock: func(id int32) {
				mock.ExpectExec(q).WithArgs(id).WillReturnResult(pgxmock.NewResult("UPDATED", 1))
			},
		},
		{
			name:    "exec error",
			wantErr: true,
			mock: func(id int32) {
				mock.ExpectExec(q).WithArgs(id).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.id)

			err := repo.DisableProject(context.Background(), tc.id)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestProjectRepository_GetParticipant(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		SELECT 
		    pp.user_id, pp.project_id, pp.role, pp.position, pp.created_at, pp.updated_at,
		    u.name, u.surname, u.username
		FROM project_participants pp
			JOIN users u ON u.id=pp.user_id
		WHERE pp.user_id=$1 AND pp.project_id=$2`

	tempTime := time.Now()

	tests := []struct {
		name        string
		response    *domain.ProjectParticipant
		uid, pid    int32
		wantErr     bool
		expectedErr error
		mock        func(uid, pid int32)
	}{
		{
			name: "should pass",
			uid:  1, pid: 1,
			response: &domain.ProjectParticipant{
				UserID:    1,
				ProjectID: 1,
				Name:      "name",
				Surname:   "surname",
				Username:  "username",
				Role:      domain.UserRoleUser,
				Position:  domain.UserPositionBackend,
				CreatedAt: tempTime,
				UpdatedAt: tempTime,
			},
			wantErr: false,
			mock: func(uid, pid int32) {
				row := mock.NewRows([]string{
					"pp.user_id",
					"pp.project_id",
					"pp.role",
					"pp.position",
					"pp.created_at",
					"pp.updated_at",
					"u.name",
					"u.surname",
					"u.username",
				}).
					AddRow(
						int32(1),
						int32(1),
						domain.UserRoleUser,
						domain.UserPositionBackend,
						tempTime,
						tempTime,
						"name",
						"surname",
						"username",
					)

				mock.ExpectQuery(q).WithArgs(uid, pid).WillReturnRows(row)
			},
		},
		{
			name: "query error",
			uid:  1, pid: 1,
			wantErr: true,
			mock: func(uid, pid int32) {
				mock.ExpectQuery(q).WithArgs(uid, pid).WillReturnError(errors.New("some error"))
			},
		},
		{
			name: "no rows error",
			uid:  1, pid: 1,
			wantErr:     true,
			expectedErr: repository.ErrObjectNotFound,
			mock: func(uid, pid int32) {
				mock.ExpectQuery(q).WithArgs(uid, pid).WillReturnError(pgx.ErrNoRows)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.uid, tc.pid)

			resp, err := repo.GetParticipant(context.Background(), tc.uid, tc.pid)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)

				if tc.expectedErr != nil {
					require.ErrorIs(tt, err, tc.expectedErr)
				}
				return
			}

			require.NoError(tt, err)
			require.Equal(tt, tc.response, resp)
		})
	}
}

func TestProjectRepository_GetParticipants(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		SELECT
		    pp.user_id, pp.project_id, pp.role, pp.position, pp.created_at, pp.updated_at,
		    u.name, u.surname, u.username
	 	FROM project_participants pp
			JOIN users u ON u.id = pp.user_id
	 	WHERE pp.project_id = $1`

	tempTime := time.Now()

	tests := []struct {
		name     string
		pid      int32
		response []domain.ProjectParticipant
		wantErr  bool
		mock     func(pid int32)
	}{
		{
			name: "should pass",
			pid:  1,
			response: []domain.ProjectParticipant{
				{
					UserID:    1,
					ProjectID: 1,
					Name:      "name1",
					Surname:   "surname1",
					Username:  "username1",
					Role:      1,
					Position:  1,
					CreatedAt: tempTime,
					UpdatedAt: tempTime,
				},
				{
					UserID:    2,
					ProjectID: 1,
					Name:      "name2",
					Surname:   "surname2",
					Username:  "username2",
					Role:      2,
					Position:  2,
					CreatedAt: tempTime,
					UpdatedAt: tempTime,
				},
			},
			wantErr: false,
			mock: func(pid int32) {
				rows := mock.NewRows([]string{
					"pp.user_id",
					"pp.project_id",
					"pp.role",
					"pp.position",
					"pp.created_at",
					"pp.updated_at",
					"u.name",
					"u.surname",
					"u.username",
				}).
					AddRow(
						int32(1),
						pid,
						domain.UserRole(1),
						domain.UserPosition(1),
						tempTime,
						tempTime,
						"name1",
						"surname1",
						"username1",
					).
					AddRow(
						int32(2),
						pid,
						domain.UserRole(2),
						domain.UserPosition(2),
						tempTime,
						tempTime,
						"name2",
						"surname2",
						"username2",
					)

				mock.ExpectQuery(q).
					WithArgs(pid).
					WillReturnRows(rows).
					RowsWillBeClosed()
			},
		},
		{
			name:     "scan error",
			response: nil,
			wantErr:  true,
			mock: func(pid int32) {
				mock.ExpectQuery(q).WithArgs(pid).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.pid)

			resp, err := repo.GetParticipants(context.Background(), tc.pid)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
			require.Equal(tt, tc.response, resp)
		})
	}
}

func TestProjectRepository_AddParticipant(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		INSERT INTO project_participants(project_id, user_id, role, position) 
		VALUES ($1,$2,$3,$4)`

	tests := []struct {
		name     string
		response int32
		wantErr  bool
		mock     func(pp *domain.ProjectParticipant)
	}{
		{
			name:     "should pass",
			response: 1,
			wantErr:  false,
			mock: func(pp *domain.ProjectParticipant) {
				mock.
					ExpectExec(q).
					WithArgs(
						pp.ProjectID,
						pp.UserID,
						pp.Role,
						pp.Position,
					).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
		},
		{
			name:     "query error",
			response: 0,
			wantErr:  true,
			mock: func(pp *domain.ProjectParticipant) {
				mock.
					ExpectExec(q).
					WithArgs(
						pp.ProjectID,
						pp.UserID,
						pp.Role,
						pp.Position,
					).WillReturnError(fmt.Errorf("some error"))
			},
		},
	}

	testProjectParticipant := &domain.ProjectParticipant{
		UserID:    1,
		ProjectID: 2,
		Role:      3,
		Position:  4,
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(testProjectParticipant)

			err := repo.AddParticipant(context.Background(), testProjectParticipant)

			require.NoError(tt, mock.ExpectationsWereMet())

			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestProjectRepository_RemoveParticipant(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		DELETE FROM project_participants
		WHERE user_id=$1 AND project_id=$2`

	tests := []struct {
		name     string
		wantErr  bool
		uid, pid int32
		mock     func(uid, pid int32)
	}{
		{
			name: "should pass",
			uid:  1, pid: 2,
			wantErr: false,
			mock: func(uid, pid int32) {
				mock.
					ExpectExec(q).
					WithArgs(
						uid, pid,
					).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
		{
			name: "query error",
			uid:  1, pid: 2,
			wantErr: true,
			mock: func(uid, pid int32) {
				mock.
					ExpectExec(q).
					WithArgs(
						uid, pid,
					).WillReturnError(fmt.Errorf("some error"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.uid, tc.pid)

			err := repo.RemoveParticipant(context.Background(), tc.uid, tc.pid)

			require.NoError(tt, mock.ExpectationsWereMet())

			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestProjectRepository_UpdateParticipant(t *testing.T) {
	mock, repo := prepareProjectMock(t)

	q := `
		UPDATE project_participants
		SET role=$3, position=$4, updated_at=now()
		WHERE user_id=$1 AND project_id=$2`

	tests := []struct {
		name     string
		response int32
		wantErr  bool
		mock     func(pp *domain.ProjectParticipant)
	}{
		{
			name:     "should pass",
			response: 1,
			wantErr:  false,
			mock: func(pp *domain.ProjectParticipant) {
				mock.
					ExpectExec(q).
					WithArgs(
						pp.UserID,
						pp.ProjectID,
						pp.Role,
						pp.Position,
					).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
		},
		{
			name:     "query error",
			response: 0,
			wantErr:  true,
			mock: func(pp *domain.ProjectParticipant) {
				mock.
					ExpectExec(q).
					WithArgs(
						pp.UserID,
						pp.ProjectID,
						pp.Role,
						pp.Position,
					).WillReturnError(fmt.Errorf("some error"))
			},
		},
	}

	testProjectParticipant := &domain.ProjectParticipant{
		UserID:    1,
		ProjectID: 2,
		Role:      3,
		Position:  4,
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(testProjectParticipant)

			err := repo.UpdateParticipant(context.Background(), testProjectParticipant)

			require.NoError(tt, mock.ExpectationsWereMet())

			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}
