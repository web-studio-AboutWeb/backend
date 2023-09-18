package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/handler/http/dto"
	smocks "web-studio-backend/internal/app/handler/http/mocks"
)

func user(t *testing.T) (*userHandler, *smocks.MockUserService) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	userService := smocks.NewMockUserService(mockCtl)
	userHandler := newUserHandler(userService)

	return userHandler, userService
}

func TestUserHandler_CreateUser(t *testing.T) {
	handler, serv := user(t)

	ctx := context.Background()

	type test struct {
		name     string
		response *domain.User
		err      *apperror.Error
		wantErr  bool
	}

	tests := []struct {
		test
		in   *dto.CreateUserIn
		code int
		mock func(user *dto.CreateUserIn)
	}{
		{
			test: test{
				name:    "should pass",
				err:     nil,
				wantErr: false,
				response: &domain.User{
					ID:           1,
					Name:         "name",
					Surname:      "surname",
					Login:        "login",
					Password:     "",
					Role:         1,
					RoleName:     "User",
					Position:     1,
					PositionName: "Frontend",
				},
			},
			in: &dto.CreateUserIn{
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "password",
				Role:     1,
				Position: 1,
			},
			code: http.StatusOK,
			mock: func(user *dto.CreateUserIn) {
				serv.EXPECT().CreateUser(ctx, user.ToDomain()).Return(&domain.User{
					ID:           1,
					Name:         user.Name,
					Surname:      user.Surname,
					Login:        user.Login,
					Password:     "",
					Role:         user.Role,
					RoleName:     user.Role.String(),
					Position:     user.Position,
					PositionName: user.Position.String(),
				}, nil)
			},
		},
		{
			test: test{
				name:     "internal error",
				wantErr:  true,
				response: nil,
			},
			in:   &dto.CreateUserIn{},
			code: http.StatusInternalServerError,
			mock: func(user *dto.CreateUserIn) {
				serv.EXPECT().CreateUser(ctx, user.ToDomain()).Return(nil, fmt.Errorf("some error"))
			},
		},
		{
			test: test{
				name:    "login already taken error",
				wantErr: true,
				err: &apperror.Error{
					Field: "login",
					Type:  apperror.InvalidRequestType,
				},
				response: nil,
			},
			in:   &dto.CreateUserIn{},
			code: http.StatusBadRequest,
			mock: func(user *dto.CreateUserIn) {
				serv.EXPECT().CreateUser(ctx, user.ToDomain()).Return(nil, &apperror.Error{
					Message: "Login already taken.",
					Field:   "login",
					Type:    apperror.InvalidRequestType,
				})
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.in)

			rec := httptest.NewRecorder()

			body := new(bytes.Buffer)

			data, err := json.Marshal(tc.in)
			require.NoError(tt, err)

			_, err = body.Write(data)
			require.NoError(tt, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", body)

			handler.createUser(rec, req)

			result := rec.Result()
			defer result.Body.Close()

			require.Equal(tt, tc.code, result.StatusCode)

			data, err = io.ReadAll(result.Body)
			require.NoError(tt, err)

			if tc.wantErr {
				require.NotEqual(tt, http.StatusOK, result.StatusCode)

				if tc.err != nil {
					var out *apperror.Error
					err = json.Unmarshal(data, &out)
					require.NoError(tt, err)

					require.Equal(tt, tc.err.Field, out.Field)
					require.Equal(tt, tc.err.Type, out.Type)
				}

				return
			}

			var out *domain.User
			err = json.Unmarshal(data, &out)
			require.NoError(tt, err)

			require.Equal(tt, out, tc.response)
		})
	}
}
