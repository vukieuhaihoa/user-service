package user

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// UpdateUserByID updates an existing user in the database by their ID.
// It takes a context, an ID, and a user model with updated details as input.
// Returns an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be updated.
//   - updatedUser: The user model containing the updated details.
//
// Returns:
//   - error: An error if the update fails, otherwise nil.
func (u *userRepository) UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error {
	result := u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updatedUser)
	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return dbutils.ErrRecordNotFoundType
	}

	return nil
}
