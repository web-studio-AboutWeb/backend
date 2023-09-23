package service_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/service"
	"web-studio-backend/internal/app/service/mocks"
)

func user(t *testing.T) (*service.UserService, *mocks.MockUserRepository, *mocks.MockFileRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtl)
	fileRepo := mocks.NewMockFileRepository(mockCtl)
	userService := service.NewUserService(userRepo, fileRepo)

	return userService, userRepo, fileRepo
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	serv, repo, _ := user(t)

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
				repo.EXPECT().CheckUserUniqueness(ctx, user.Username, user.Email).Return(nil, nil)
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
				repo.EXPECT().CheckUserUniqueness(ctx, user.Username, user.Email).Return(&domain.User{Username: "123"}, nil)
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
				repo.EXPECT().CheckUserUniqueness(ctx, user.Username, user.Email).Return(nil, nil)
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
				repo.EXPECT().CheckUserUniqueness(ctx, user.Username, user.Email).Return(nil, nil)
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

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()

	serv, repo, _ := user(t)

	ctx := context.Background()

	tempTime := time.Now()

	type test struct {
		name    string
		res     *domain.User
		err     error
		wantErr bool
	}

	tests := []struct {
		test
		id   int32
		mock func(id int32)
	}{
		{
			test: test{
				name: "should pass",
				res: &domain.User{
					ID:         1,
					Name:       "name",
					Surname:    "surname",
					Username:   "login",
					Email:      "email",
					Role:       2,
					IsTeamLead: false,
					CreatedAt:  tempTime,
					UpdatedAt:  tempTime,
					DisabledAt: &tempTime,
					ImageID:    "image_id",
				},
				wantErr: false,
			},
			id: 1,
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(&domain.User{
					ID:         1,
					Name:       "name",
					Surname:    "surname",
					Username:   "login",
					Email:      "email",
					Role:       2,
					IsTeamLead: false,
					CreatedAt:  tempTime,
					UpdatedAt:  tempTime,
					DisabledAt: &tempTime,
					ImageID:    "image_id",
				}, nil)
			},
		},
		{
			test: test{
				name:    "error",
				wantErr: true,
			},
			id: 1,
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(nil, errors.New("unknown"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id)

			res, err := serv.GetUser(ctx, tc.id)
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

func TestUserService_GetUsers(t *testing.T) {
	t.Parallel()

	serv, repo, _ := user(t)

	ctx := context.Background()

	tempTime := time.Now()

	type test struct {
		name    string
		res     []domain.User
		err     error
		wantErr bool
	}

	tests := []struct {
		test
		mock func()
	}{
		{
			test: test{
				name: "should pass",
				res: []domain.User{
					{
						ID:         1,
						Name:       "name1",
						Surname:    "surname1",
						Username:   "login1",
						Email:      "email1",
						Role:       1,
						IsTeamLead: false,
						CreatedAt:  tempTime,
						UpdatedAt:  tempTime,
						DisabledAt: &tempTime,
						ImageID:    "image_id1",
					},
					{
						ID:         2,
						Name:       "name2",
						Surname:    "surname2",
						Username:   "login2",
						Email:      "email2",
						Role:       2,
						IsTeamLead: true,
						CreatedAt:  tempTime,
						UpdatedAt:  tempTime,
						DisabledAt: &tempTime,
						ImageID:    "image_id2",
					},
				},
			},
			mock: func() {
				repo.EXPECT().GetUsers(ctx).Return([]domain.User{
					{
						ID:         1,
						Name:       "name1",
						Surname:    "surname1",
						Username:   "login1",
						Email:      "email1",
						Role:       1,
						IsTeamLead: false,
						CreatedAt:  tempTime,
						UpdatedAt:  tempTime,
						DisabledAt: &tempTime,
						ImageID:    "image_id1",
					},
					{
						ID:         2,
						Name:       "name2",
						Surname:    "surname2",
						Username:   "login2",
						Email:      "email2",
						Role:       2,
						IsTeamLead: true,
						CreatedAt:  tempTime,
						UpdatedAt:  tempTime,
						DisabledAt: &tempTime,
						ImageID:    "image_id2",
					},
				}, nil)
			},
		},
		{
			test: test{
				name:    "error",
				wantErr: true,
			},
			mock: func() {
				repo.EXPECT().GetUsers(ctx).Return(nil, errors.New("unknown"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, err := serv.GetUsers(ctx)
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

func TestUserService_UpdateUser(t *testing.T) {
	t.Parallel()

	serv, repo, _ := user(t)

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
					ID:      1,
					Name:    "name",
					Surname: "surname",
					Role:    1,
				},
				err:     nil,
				wantErr: false,
			},
			user: &domain.User{
				ID:      1,
				Name:    "name",
				Surname: "surname",
				Role:    1,
			},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUser(ctx, user.ID).Return(&domain.User{
					ID:      1,
					Name:    "name",
					Surname: "surname",
					Role:    1,
				}, nil)
				repo.EXPECT().UpdateUser(ctx, user).Return(nil)
				repo.EXPECT().GetUser(ctx, user.ID).Return(&domain.User{
					ID:      1,
					Name:    "name",
					Surname: "surname",
					Role:    1,
				}, nil)
			},
		},
		{
			test: test{
				name:    "update user error",
				res:     nil,
				wantErr: true,
			},
			user: &domain.User{},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUser(ctx, user.ID).Return(&domain.User{
					ID:      1,
					Name:    "name",
					Surname: "surname",
					Role:    1,
				}, nil)
				repo.EXPECT().UpdateUser(ctx, user).Return(errors.New("update user error"))
			},
		},
		{
			test: test{
				name:    "get user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUser(ctx, user.ID).Return(nil, errors.New("get user error"))
			},
		},
		{
			test: test{
				name:    "get updated user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			user: &domain.User{},
			mock: func(user *domain.User) {
				repo.EXPECT().GetUser(ctx, user.ID).Return(&domain.User{
					ID:      1,
					Name:    "name",
					Surname: "surname",
					Role:    1,
				}, nil)
				repo.EXPECT().UpdateUser(ctx, user).Return(nil)
				repo.EXPECT().GetUser(ctx, user.ID).Return(nil, errors.New("get user error"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.user)

			res, err := serv.UpdateUser(ctx, tc.user)
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

func TestUserService_GetUserImage(t *testing.T) {
	t.Parallel()

	serv, repo, fileRepo := user(t)

	ctx := context.Background()

	type test struct {
		name    string
		res     *domain.User
		err     error
		wantErr bool
	}

	tests := []struct {
		test
		id   int32
		mock func(id int32)
	}{
		{
			test: test{
				name: "should pass",
				res: &domain.User{
					ID:           1,
					ImageID:      "image_id",
					ImageContent: []byte("image_content"),
				},
				err:     nil,
				wantErr: false,
			},
			id: 1,
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(&domain.User{
					ID:      1,
					ImageID: "image_id",
				}, nil)
				fileRepo.EXPECT().Read(ctx, filepath.Join("users", "image_id")).Return([]byte("image_content"), nil)
			},
		},
		{
			test: test{
				name:    "get user error",
				res:     nil,
				err:     nil,
				wantErr: true,
			},
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(nil, errors.New("get user error"))
			},
		},
		{
			test: test{
				name:    "not found",
				res:     nil,
				wantErr: true,
			},
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(&domain.User{
					ID:      id,
					ImageID: "",
				}, nil)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id)

			res, err := serv.GetUserImage(ctx, tc.id)
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

func TestUserService_RemoveUser(t *testing.T) {
	t.Parallel()

	serv, repo, _ := user(t)

	ctx := context.Background()

	type test struct {
		name    string
		err     error
		wantErr bool
	}

	tests := []struct {
		test
		id   int32
		mock func(id int32)
	}{
		{
			test: test{
				name:    "should pass",
				err:     nil,
				wantErr: false,
			},
			id: 1,
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(&domain.User{
					ID: 1,
				}, nil)
				repo.EXPECT().DisableUser(ctx, id).Return(nil)
			},
		},
		{
			test: test{
				name:    "get user error",
				err:     nil,
				wantErr: true,
			},
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(nil, errors.New("get user error"))
			},
		},
		{
			test: test{
				name:    "disable error",
				wantErr: true,
			},
			mock: func(id int32) {
				repo.EXPECT().GetUser(ctx, id).Return(&domain.User{
					ID: id,
				}, nil)
				repo.EXPECT().DisableUser(ctx, id).Return(errors.New("some err"))
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.id)

			err := serv.RemoveUser(ctx, tc.id)
			if !tc.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

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
