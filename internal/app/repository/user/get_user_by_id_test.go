package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
	"gorm.io/gorm"
)

func TestUser_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB
		inputID string

		expectedError  error
		expectedOutput *model.User
	}{
		{
			name: "Get user by ID successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "de305d54-75b4-431b-adb2-eb6b9e546000",

			expectedOutput: &model.User{
				Base: model.Base{
					ID:        "de305d54-75b4-431b-adb2-eb6b9e546000",
					CreatedAt: fixture.TestTime,
					UpdatedAt: fixture.TestTime,
				},
				Username:    "Alice",
				DisplayName: "Alice",
				Email:       "alice@example.com",
			},
		},
		{
			name: "Get user by ID failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "non-existent-id",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUserRepository(db)

			res, err := testUserRepo.GetUserByID(ctx, tc.inputID)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			res.Password = "" // omit password field for comparison
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
