package user

import (
	"context"

	"github.com/vukieuhaihoa/user-service/internal/app/model"
)

// CreateUser creates a new user with the provided information.
// It hashes the password before storing the user in the database.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - username: The username of the new user.
//   - password: The password of the new user.
//   - displayName: The display name of the new user.
//   - email: The email address of the new user.
//
// Returns:
//   - *model.User: The created user model.
//   - error: An error if the creation fails, otherwise nil.
func (u *userService) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashedPassword, err := u.passwordHashing.Hash(password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Username:    username,
		Password:    hashedPassword,
		DisplayName: displayName,
		Email:       email,
	}

	createdUser, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
