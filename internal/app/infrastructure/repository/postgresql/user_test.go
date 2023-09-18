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
		response int16
		wantErr  bool
		mock     func(user *domain.User, id int16)
	}{
		{
			name: "should pass",
			in: &domain.User{
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "password",
				Role:     1,
				Position: 1,
			},
			response: 1,
			wantErr:  false,
			mock: func(user *domain.User, id int16) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(id)

				mock.
					ExpectQuery(q).
					WithArgs(
						user.Name,
						user.Surname,
						user.Login,
						user.Password,
						user.Role,
						user.Position,
					).WillReturnRows(rows)
			},
		},
		{
			name: "query error",
			in: &domain.User{
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "password",
				Role:     1,
				Position: 1,
			},
			response: 0,
			wantErr:  true,
			mock: func(user *domain.User, id int16) {
				mock.
					ExpectQuery(q).
					WithArgs(
						user.Name,
						user.Surname,
						user.Login,
						user.Password,
						user.Role,
						user.Position,
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
