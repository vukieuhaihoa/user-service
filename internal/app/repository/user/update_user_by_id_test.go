package user

import (
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
)

func TestUser_UpdateUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		inputID       string
		inputUserData *model.User

		expectedError error
	}{
		{
			name: "Update user by ID successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "de305d54-75b4-431b-adb2-eb6b9e546000",

			inputUserData: &model.User{
				DisplayName: "Alice Updated",
				Email:       "alice.updated@example.com",
			},
		},
		{
			name: "Update user by ID failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "non-existent-id",

			inputUserData: &model.User{
				DisplayName: "Non Existent User",
				Email:       "",
			},

			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Update user by ID failed - duplicate email",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputID: "de305d54-75b4-431b-adb2-eb6b9e546000",

			inputUserData: &model.User{
				DisplayName: "Alice",
				Email:       "bob@example.com", // duplicate email
			},

			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testUserRepo := NewUserRepository(db)

			err := testUserRepo.UpdateUserByID(ctx, tc.inputID, tc.inputUserData)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
		})
	}
}
