package postgresql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository/postgresql"
)

func TestUserRepository_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := postgresql.NewUserRepository(mock)

	q := `^INSERT INTO users(.+) VALUES(.+) RETURNING id$`

	tests := []struct {
		name     string
		response int32
		wantErr  bool
		mock     func(user *domain.User, id int32)
	}{
		{
			name:     "should pass",
			response: 1,
			wantErr:  false,
			mock: func(user *domain.User, id int32) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(id)

				mock.
					ExpectQuery(q).
					WithArgs(
						user.Name,
						user.Surname,
						user.Username,
						user.Email,
						user.EncodedPassword,
						user.Salt,
						user.Role,
						user.IsTeamLead,
						user.ImageID,
					).WillReturnRows(rows)
			},
		},
		{
			name:     "query error",
			response: 0,
			wantErr:  true,
			mock: func(user *domain.User, id int32) {
				mock.
					ExpectQuery(q).
					WithArgs(
						user.Name,
						user.Surname,
						user.Username,
						user.Email,
						user.EncodedPassword,
						user.Salt,
						user.Role,
						user.IsTeamLead,
						user.ImageID,
					).WillReturnError(fmt.Errorf("some error"))
			},
		},
	}

	testUser := &domain.User{
		Name:            "name",
		Surname:         "surname",
		Username:        "login",
		EncodedPassword: "password",
		Role:            1,
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(testUser, tc.response)

			resp, err := repo.CreateUser(context.Background(), testUser)

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

func TestUserRepository_GetUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := postgresql.NewUserRepository(mock)

	q := `^SELECT (.+) FROM users WHERE (.+)$`

	tempTime := time.Now()

	tests := []struct {
		name     string
		response *domain.User
		in       int32
		wantErr  bool
		mock     func(id int32)
	}{
		{
			name: "should pass",
			in:   1,
			response: &domain.User{
				ID:         1,
				Name:       "name",
				Surname:    "surname",
				Username:   "username",
				Email:      "email",
				Role:       1,
				RoleName:   "",
				IsTeamLead: true,
				CreatedAt:  tempTime,
				UpdatedAt:  tempTime,
				DisabledAt: nil,
				ImageID:    "test",
			},
			wantErr: false,
			mock: func(id int32) {
				row := mock.NewRows([]string{
					"id", "name", "surname", "username", "email",
					"created_at", "updated_at", "disabled_at", "role",
					"is_teamlead", "image_id",
				}).
					AddRow(
						id,
						"name",
						"surname",
						"username",
						"email",
						tempTime,
						tempTime,
						nil,
						domain.UserRole(1),
						true,
						"test",
					)

				mock.ExpectQuery(q).WithArgs(id).WillReturnRows(row)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in)

			resp, err := repo.GetUser(context.Background(), tc.in)

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
