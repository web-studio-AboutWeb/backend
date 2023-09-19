package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/service"
	"web-studio-backend/internal/app/service/mocks"
)

func user(t *testing.T) (*service.UserService, *mocks.MockUserRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtl)
	userService := service.NewUserService(userRepo)

	return userService, userRepo
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	serv, repo := user(t)

	ctx := context.Background()

	type test struct {
		name    string
		res     *domain.User
		err     error
		wantErr bool
	}

	tests := []struct {
		test
		user *domain.User
		mock func(user *domain.User)
	}{
		{
			test: test{
				name: "should pass",
				res: &domain.User{
					ID:       1,
					Name:     "name",
					Surname:  "surname",
					Username: "login",
					Email:    "email@mail.com",
					Role:     1,
				},
				err:     nil,
				wantErr: false,
			},
			user: &domain.User{
				ID:              1,
				Name:            "name",
				Surname:         "surname",
				Username:        "login",
				Email:           "email@mail.com",
				EncodedPassword: "123",
				Role:            1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().CheckUsernameUniqueness(ctx, user.Username, user.Email).Return(nil, nil)
				repo.EXPECT().CreateUser(ctx, user).Return(int32(1), nil)
				repo.EXPECT().GetUser(ctx, int32(1)).Return(&domain.User{
					ID:       1,
					Name:     user.Name,
					Surname:  user.Surname,
					Username: user.Username,
					Email:    user.Email,
					Role:     user.Role,
				}, nil)
			},
		},
		{
			test: test{
				name:    "login already taken error",
				res:     nil,
				wantErr: true,
			},
			user: &domain.User{
				ID:              1,
				Name:            "name",
				Surname:         "surname",
				Username:        "login",
				Email:           "email@mail.com",
				EncodedPassword: "123",
				Role:            1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().CheckUsernameUniqueness(ctx, user.Username, user.Email).Return(&domain.User{Username: "123"}, nil)
			},
		},
		{
			test: test{
				name:    "create user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{Username: "123"},
			mock: func(user *domain.User) {
				repo.EXPECT().CheckUsernameUniqueness(ctx, user.Username, user.Email).Return(nil, nil)
				repo.EXPECT().CreateUser(ctx, user).Return(int32(0), fmt.Errorf("create user error"))
			},
		},
		{
			test: test{
				name:    "get user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{Username: "login"},
			mock: func(user *domain.User) {
				repo.EXPECT().CheckUsernameUniqueness(ctx, user.Username, user.Email).Return(nil, nil)
				repo.EXPECT().CreateUser(ctx, user).Return(int32(1), nil)
				repo.EXPECT().GetUser(ctx, int32(1)).Return(nil, fmt.Errorf("get user error"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.user)

			res, err := serv.CreateUser(ctx, tc.user)
			if !tc.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Equal(t, tc.res, res)

			if !tc.wantErr {
				return
			}

			require.Error(t, err)
			if tc.err != nil {
				require.ErrorAs(t, err, tc.err)
			}
		})
	}
}
