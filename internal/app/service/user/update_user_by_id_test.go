package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	mockUserRepo "github.com/vukieuhaihoa/user-service/internal/app/repository/user/mocks"
)

func TestService_UpdateUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo func(ctx context.Context) *mockUserRepo.Repository
		inputUserID       string
		inputDisplayName  string
		inputEmail        string

		expectedError error
	}{
		{
			name: "Update user by ID successfully",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", &model.User{
					DisplayName: "Updated User",
					Email:       "updateduser@example.com",
				}).Return(nil)
				return repoMock
			},

			inputUserID:      "de305d54-75b4-431b-adb2-eb6b9e546099",
			inputDisplayName: "Updated User",
			inputEmail:       "updateduser@example.com",
		},
		{
			name: "Fail to update user by ID - user not found",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("UpdateUserByID", ctx, "nonexistentid", &model.User{
					DisplayName: "Updated User",
					Email:       "updateduser@example.com",
				}).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputUserID:      "nonexistentid",
			inputDisplayName: "Updated User",
			inputEmail:       "updateduser@example.com",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Fail to update user by ID - duplicate email",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", &model.User{
					DisplayName: "Updated User",
					Email:       "duplicateemail@example.com",
				}).Return(dbutils.ErrDuplicationType)
				return repoMock
			},

			inputUserID:      "de305d54-75b4-431b-adb2-eb6b9e546099",
			inputDisplayName: "Updated User",
			inputEmail:       "duplicateemail@example.com",

			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			userRepoMock := tc.setupMockUserRepo(ctx)

			userService := NewUserService(userRepoMock, nil, nil)

			err := userService.UpdateUserByID(ctx, tc.inputUserID, tc.inputDisplayName, tc.inputEmail)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
