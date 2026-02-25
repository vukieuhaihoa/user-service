package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
	"gorm.io/gorm"
)

func TestUser_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		inputUsername string

		expectedError  error
		expectedOutput *model.User
	}{
		{
			name: "Get user by username successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUsername: "Bob",

			expectedOutput: &model.User{
				Base: model.Base{
					ID:        "123e4567-e89b-12d3-a456-eb6b9e546001",
					CreatedAt: fixture.TestTime,
					UpdatedAt: fixture.TestTime,
				},
				Username:    "Bob",
				DisplayName: "Bob",
				Email:       "bob@example.com",
			},
		},
		{
			name: "Get user by username failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUsername: "NonExistentUser",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUserRepository(db)

			res, err := testUserRepo.GetUserByUsername(ctx, tc.inputUsername)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			res.Password = "" // omit password field for comparison
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
