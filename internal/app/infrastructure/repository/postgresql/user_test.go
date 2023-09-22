package postgresql_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository/postgresql"
)

func prepareUserMock(t *testing.T) (pgxmock.PgxPoolIface, *postgresql.UserRepository) {
	t.Helper()

	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	repo := postgresql.NewUserRepository(mock)

	return mock, repo
}

func TestUserRepository_CreateUser(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `
		INSERT INTO users(name, surname, username, email, encoded_password, salt, role, is_teamlead, image_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING  id`

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
		Email:           "email",
		EncodedPassword: "password",
		Salt:            "salt",
		Role:            1,
		ImageID:         "image_id",
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
	mock, repo := prepareUserMock(t)

	q := `SELECT 
		    id, name, surname, username, email, created_at, updated_at, disabled_at, role, is_teamlead, image_id
        FROM users
        WHERE id = $1`

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
		{
			name:    "query error",
			in:      1,
			wantErr: true,
			mock: func(id int32) {
				mock.ExpectQuery(q).WithArgs(id).WillReturnError(errors.New("some error"))
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

func TestUserRepository_GetActiveUser(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `
		SELECT
		    id, name, surname, username, email, created_at, 
		    updated_at, disabled_at, role, is_teamlead, image_id
        FROM users
        WHERE id = $1 AND disabled_at IS NULL`

	tempTime := time.Now()

	tests := []struct {
		name     string
		in       int32
		response *domain.User
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
		{
			name:     "scan error",
			in:       1,
			response: nil,
			wantErr:  true,
			mock: func(id int32) {
				mock.ExpectQuery(q).WithArgs(id).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in)

			resp, err := repo.GetActiveUser(context.Background(), tc.in)

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

func TestUserRepository_CheckUserUniqueness(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `
		SELECT id, username, email
        FROM users
        WHERE (lower(username)=lower($1) OR lower(email)=lower($2))
          		AND disabled_at IS NULL`

	tests := []struct {
		name            string
		username, email string
		response        *domain.User
		wantErr         bool
		mock            func(username, email string)
	}{
		{
			name:     "should pass",
			username: "test",
			email:    "test@test.com",
			response: &domain.User{
				ID:       1,
				Username: "test",
				Email:    "test@test.com",
			},
			wantErr: false,
			mock: func(username, email string) {
				row := mock.NewRows([]string{
					"id", "username", "email",
				}).
					AddRow(
						int32(1),
						username,
						email,
					)

				mock.ExpectQuery(q).WithArgs(username, email).WillReturnRows(row)
			},
		},
		{
			name:     "scan error",
			response: nil,
			wantErr:  true,
			mock: func(username, email string) {
				mock.ExpectQuery(q).WithArgs(username, email).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.username, tc.email)

			resp, err := repo.CheckUserUniqueness(context.Background(), tc.username, tc.email)

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

func TestUserRepository_GetUsers(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `
		SELECT 
		    id, name, surname, username, email, created_at, updated_at, disabled_at,
		    role, is_teamlead, image_id
        FROM users`

	tempTime := time.Now()

	tests := []struct {
		name     string
		response []domain.User
		wantErr  bool
		mock     func()
	}{
		{
			name: "should pass",
			response: []domain.User{
				{
					ID:         1,
					Name:       "name1",
					Surname:    "surname1",
					Username:   "username1",
					Email:      "email1",
					Role:       domain.UserRole(1),
					IsTeamLead: true,
					CreatedAt:  tempTime,
					UpdatedAt:  tempTime,
					DisabledAt: nil,
					ImageID:    "image1",
				},
				{
					ID:         2,
					Name:       "name2",
					Surname:    "surname2",
					Username:   "username2",
					Email:      "email2",
					Role:       domain.UserRole(2),
					IsTeamLead: false,
					CreatedAt:  tempTime,
					UpdatedAt:  tempTime,
					DisabledAt: &tempTime,
					ImageID:    "image2",
				},
			},
			wantErr: false,
			mock: func() {
				rows := mock.NewRows([]string{
					"id", "name", "surname", "username", "email", "created_at", "updated_at", "disabled_at",
					"role", "is_teamlead", "image_id",
				}).
					AddRow(
						int32(1),
						"name1",
						"surname1",
						"username1",
						"email1",
						tempTime,
						tempTime,
						nil,
						domain.UserRole(1),
						true,
						"image1",
					).
					AddRow(
						int32(2),
						"name2",
						"surname2",
						"username2",
						"email2",
						tempTime,
						tempTime,
						&tempTime,
						domain.UserRole(2),
						false,
						"image2",
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

			resp, err := repo.GetUsers(context.Background())

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

func TestUserRepository_DisableUser(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `UPDATE users SET disabled_at=now() WHERE id=$1`

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

			err := repo.DisableUser(context.Background(), tc.id)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `
		UPDATE users 
		SET name=$2, surname=$3, role=$4, updated_at=now()
		WHERE id = $1`

	tests := []struct {
		name    string
		in      *domain.User
		wantErr bool
		mock    func(user *domain.User)
	}{
		{
			name:    "should pass",
			wantErr: false,
			in: &domain.User{
				ID:      1,
				Name:    "name",
				Surname: "surname",
				Role:    domain.UserRoleUser,
			},
			mock: func(user *domain.User) {
				mock.ExpectExec(q).WithArgs(
					user.ID,
					user.Name,
					user.Surname,
					user.Role,
				).WillReturnResult(pgxmock.NewResult("UPDATED", 1))
			},
		},
		{
			name:    "exec error",
			wantErr: true,
			in:      &domain.User{},
			mock: func(user *domain.User) {
				mock.ExpectExec(q).WithArgs(
					user.ID,
					user.Name,
					user.Surname,
					user.Role,
				).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in)

			err := repo.UpdateUser(context.Background(), tc.in)

			require.NoError(tt, mock.ExpectationsWereMet())
			if tc.wantErr {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestUserRepository_GetUserByLogin(t *testing.T) {
	mock, repo := prepareUserMock(t)

	q := `
		SELECT 
		    id, username, email, encoded_password, salt
        FROM users
        WHERE (lower(username)=lower($1) OR lower(email)=lower($1))
          		AND disabled_at IS NULL`

	tests := []struct {
		name     string
		login    string
		response *domain.User
		wantErr  bool
		mock     func(login string)
	}{
		{
			name:  "should pass",
			login: "login1",
			response: &domain.User{
				ID:              1,
				Username:        "username",
				Email:           "email",
				EncodedPassword: "encoded_password",
				Salt:            "salt",
			},
			wantErr: false,
			mock: func(login string) {
				row := mock.NewRows([]string{
					"id", "username", "email", "encoded_password", "salt",
				}).
					AddRow(
						int32(1),
						"username",
						"email",
						"encoded_password",
						"salt",
					)

				mock.ExpectQuery(q).WithArgs(login).WillReturnRows(row)
			},
		},
		{
			name:     "scan error",
			response: nil,
			wantErr:  true,
			mock: func(login string) {
				mock.ExpectQuery(q).WithArgs(login).WillReturnError(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.login)

			resp, err := repo.GetUserByLogin(context.Background(), tc.login)

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
