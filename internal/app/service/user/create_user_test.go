package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	mockPasswordHashing "github.com/vukieuhaihoa/bookmark-libs/pkg/utils/mocks"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	mockUserRepo "github.com/vukieuhaihoa/user-service/internal/app/repository/user/mocks"
)

func TestService_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockPasswordHashing func(t *testing.T) *mockPasswordHashing.PasswordHashing
		setupMockUserRepo        func(ctx context.Context) *mockUserRepo.Repository

		inputUsername    string
		inputPassword    string
		inputDisplayName string
		inputEmail       string

		expectedOutput *model.User
		expectedError  error
	}{
		{
			name: "Create user successfully",

			setupMockPasswordHashing: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", nil)
				return hashingMock
			},

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("CreateUser", ctx, &model.User{
					Username:    "testuser",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "Test User",
					Email:       "testuser@example.com",
				}).Return(&model.User{
					Base: model.Base{
						ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
					},
					Username:    "testuser",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "Test User",
					Email:       "testuser@example.com",
				}, nil)
				return repoMock
			},

			inputUsername:    "testuser",
			inputPassword:    "password123",
			inputDisplayName: "Test User",
			inputEmail:       "testuser@example.com",

			expectedOutput: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
				},
				Username:    "testuser",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				DisplayName: "Test User",
				Email:       "testuser@example.com",
			},
		},

		{
			name: "Fail to hash password",

			setupMockPasswordHashing: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("Hash", "badpassword").Return("", utils.ErrCannotGenerateHash)
				return hashingMock
			},
			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				return mockUserRepo.NewRepository(t)
			},

			inputUsername:    "testuser2",
			inputPassword:    "badpassword",
			inputDisplayName: "Test User 2",
			inputEmail:       "testuser2@example.com",

			expectedError: utils.ErrCannotGenerateHash,
		},

		{
			name: "Fail to create user in repository",

			setupMockPasswordHashing: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", nil)
				return hashingMock
			},

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("CreateUser", ctx, &model.User{
					Username:    "testuser3",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "Test User 3",
					Email:       "testuser3@example.com",
				}).Return(nil, assert.AnError)
				return repoMock
			},

			inputUsername:    "testuser3",
			inputPassword:    "password123",
			inputDisplayName: "Test User 3",
			inputEmail:       "testuser3@example.com",

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			passwordHashingMock := tc.setupMockPasswordHashing(t)
			userRepoMock := tc.setupMockUserRepo(ctx)

			userService := NewUserService(userRepoMock, passwordHashingMock, nil)

			res, err := userService.CreateUser(ctx, tc.inputUsername, tc.inputPassword, tc.inputDisplayName, tc.inputEmail)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
