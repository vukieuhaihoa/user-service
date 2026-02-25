// package user provides repository operations for user management.
// This file contains the user repository implementation using GORM for database interactions.
// It defines the User interface and its methods for creating, retrieving, and updating user records.
// The repository is essential for managing user data in the application.
package user

import (
	"context"

	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"gorm.io/gorm"
)

// Repository represents the interface for user repository operations.
//
//go:generate mockery --name=Repository --filename=user_repo.go --output=./mocks
type Repository interface {
	// CreateUser creates a new user in the database.
	// Returns the created user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - user: The user model containing the details of the user to be created.
	//
	// Returns:
	//   - *model.User: The created user model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// GetUserByUsername retrieves a user from the database by their username.
	// Returns the user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - username: The username of the user to be retrieved.
	//
	// Returns:
	//   - *model.User: The user model if found.
	//   - error: An error if the retrieval fails or the user is not found.
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetUserByID retrieves a user from the database by their ID.
	// Returns the user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be retrieved.
	//
	// Returns:
	//   - *model.User: The user model if found.
	//   - error: An error if the retrieval fails or the user is not found.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// UpdateUserByID updates an existing user in the database by their ID.
	// Returns an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be updated.
	//   - updatedUser: The user model containing the updated details.
	//
	// Returns:
	//   - error: An error if the update fails, otherwise nil.
	UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error
}

// user is the concrete implementation of the Repository interface.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of the User repository.
//
// Parameters:
//   - db: The GORM database connection.
//
// Returns:
//   - User: A new user repository instance.
func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}
