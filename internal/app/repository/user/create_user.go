package user

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// CreateUser creates a new user in the database.
// It takes a context and a user model as input and returns the created user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - newUser: The user model containing the details of the user to be created.
//
// Returns:
//   - *model.User: The created user model.
//   - error: An error if the creation fails, otherwise nil.
func (u *userRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := u.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return newUser, nil
}
