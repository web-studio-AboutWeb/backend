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
					Login:    "login",
					Password: "",
					Role:     1,
					Position: 1,
				},
				err:     nil,
				wantErr: false,
			},
			user: &domain.User{
				ID:       1,
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "",
				Role:     1,
				Position: 1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUserByLogin(ctx, user.Login).Return(nil, nil)
				repo.EXPECT().CreateUser(ctx, user).Return(int16(1), nil)
				repo.EXPECT().GetUser(ctx, int16(1)).Return(&domain.User{
					ID:       1,
					Name:     user.Name,
					Surname:  user.Surname,
					Login:    user.Login,
					Password: "",
					Role:     user.Role,
					Position: user.Position,
				}, nil)
			},
		},
		{
			test: test{
				name:    "login already taken error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{
				ID:       1,
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "",
				Role:     1,
				Position: 1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUserByLogin(ctx, user.Login).Return(&domain.User{}, nil)
			},
		},
		{
			test: test{
				name:    "create user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{
				ID:       1,
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "",
				Role:     1,
				Position: 1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUserByLogin(ctx, user.Login).Return(nil, nil)
				repo.EXPECT().CreateUser(ctx, user).Return(int16(0), fmt.Errorf("create user error"))
			},
		},
		{
			test: test{
				name:    "get user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{
				ID:       1,
				Name:     "name",
				Surname:  "surname",
				Login:    "login",
				Password: "",
				Role:     1,
				Position: 1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUserByLogin(ctx, user.Login).Return(nil, nil)
				repo.EXPECT().CreateUser(ctx, user).Return(int16(1), nil)
				repo.EXPECT().GetUser(ctx, int16(1)).Return(nil, fmt.Errorf("get user error"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.user)

			res, err := serv.CreateUser(ctx, tc.user)

			require.Equal(t, res, tc.res)

			if !tc.wantErr {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}
