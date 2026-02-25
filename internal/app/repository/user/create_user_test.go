package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"github.com/vukieuhaihoa/user-service/internal/test/fixture"
	"gorm.io/gorm"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB   func(t *testing.T) *gorm.DB
		inputUser *model.User

		expectedError  error
		expectedOutput *model.User
		verifyFunc     func(db *gorm.DB, user *model.User)
	}{
		{
			name: "Create user successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546099"},
				DisplayName: "New User",
				Username:    "New User",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "newuser@example.com",
			},

			expectedOutput: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546099"},
				DisplayName: "New User",
				Username:    "New User",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "newuser@example.com",
			},

			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}
				err := db.Where("id = ?", user.ID).First(checkUser).Error
				assert.Nil(t, err)

				// Verify timestamps are automatically set by GORM
				assert.False(t, user.CreatedAt.IsZero(), "CreatedAt should be automatically set")
				assert.False(t, user.UpdatedAt.IsZero(), "UpdatedAt should be automatically set")

				// On creation, both timestamps should be very close (within a second)
				timeDiff := user.UpdatedAt.Sub(user.CreatedAt).Abs()
				assert.Less(t, timeDiff, time.Second, "CreatedAt and UpdatedAt should be nearly equal on creation")

				// Verify other fields
				assert.Equal(t, user.ID, checkUser.ID)
				assert.Equal(t, user.Username, checkUser.Username)
				assert.Equal(t, user.Email, checkUser.Email)
				assert.Equal(t, user.DisplayName, checkUser.DisplayName)
			},
		},
		{
			name: "Create user failed - duplicate username",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546099"},
				DisplayName: "Another User",
				Username:    "Alice", // duplicate username
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "alice1@example.com",
			},

			expectedError: dbutils.ErrDuplicationType,
		},
		{
			name: "Create user failed - duplicate email",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546099"},
				DisplayName: "Another User",
				Username:    "AnotherUser",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "alice@example.com",
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

			res, err := testUserRepo.CreateUser(ctx, tc.inputUser)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			tc.verifyFunc(db, res)
		})
	}
}
