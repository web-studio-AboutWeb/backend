package postgresql_test

import (
	"context"
	"fmt"
	"testing"

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

	q := "INSERT INTO users"

	tests := []struct {
		name     string
		in       *domain.User
		response int32
		wantErr  bool
		mock     func(user *domain.User, id int32)
	}{
		{
			name: "should pass",
			in: &domain.User{
				Name:            "name",
				Surname:         "surname",
				Username:        "login",
				EncodedPassword: "password",
				Role:            1,
			},
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
			name: "query error",
			in: &domain.User{
				Name:            "name",
				Surname:         "surname",
				Username:        "login",
				EncodedPassword: "password",
				Role:            1,
			},
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

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in, tc.response)

			resp, err := repo.CreateUser(context.Background(), tc.in)

			require.NoError(tt, mock.ExpectationsWereMet())

			require.Equal(tt, tc.response, resp)

			if !tc.wantErr {
				require.NoError(tt, err)
				return
			}

			require.Error(tt, err)
		})
	}
}
