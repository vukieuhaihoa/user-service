package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	mockUserRepo "github.com/vukieuhaihoa/user-service/internal/app/repository/user/mocks"
)

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo func(ctx context.Context) *mockUserRepo.Repository

		inputUserID string

		expectedOutput *model.User
		expectedError  error
	}{
		{
			name: "Get user by ID successfully",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099").Return(&model.User{
					Base: model.Base{
						ID:        "de305d54-75b4-431b-adb2-eb6b9e546099",
						CreatedAt: fixture.TestTime,
						UpdatedAt: fixture.TestTime,
					},
					Username:    "testuser",
					DisplayName: "Test User",
					Email:       "testuser@example.com",
				}, nil)
				return repoMock
			},

			inputUserID: "de305d54-75b4-431b-adb2-eb6b9e546099",

			expectedOutput: &model.User{
				Base: model.Base{
					ID:        "de305d54-75b4-431b-adb2-eb6b9e546099",
					CreatedAt: fixture.TestTime,
					UpdatedAt: fixture.TestTime,
				},
				Username:    "testuser",
				DisplayName: "Test User",
				Email:       "testuser@example.com",
			},
		},
		{
			name: "Fail to get user by ID",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByID", ctx, "nonexistentid").Return(nil, dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputUserID:   "nonexistentid",
			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			userRepoMock := tc.setupMockUserRepo(ctx)

			userService := NewUserService(userRepoMock, nil, nil)

			res, err := userService.GetUserByID(ctx, tc.inputUserID)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
