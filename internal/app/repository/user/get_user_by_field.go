package user

import (
	"context"
	"fmt"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// GetUserByField retrieves a user from the database by a specified field and value.
// It takes a context, a field name, and a value as input and returns the user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - field: The field name to search by (e.g., "email", "username").
//   - value: The value of the field to match.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *userRepository) GetUserByField(ctx context.Context, field string, value string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return user, nil
}
