// package user provides services related to user management.
// It includes functionalities for creating users, authenticating them,
// retrieving user information, and updating user details.
package user

import (
	"context"
	"errors"
	"time"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	"github.com/vukieuhaihoa/user-service/internal/app/repository/user"
)

const TokenExpirationDuration = 24 * time.Hour

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// Service represents the interface for user service operations.
//
//go:generate mockery --name=Service --filename=user_service.go --output=./mocks
type Service interface {
	// CreateUser creates a new user with the provided information.
	// Returns the created user or an error if the operation fails.
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
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)

	// Login authenticates a user with the provided username and password.
	// Returns a JWT token if authentication is successful, or an error if it fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - username: The username of the user attempting to log in.
	//   - password: The password of the user attempting to log in.
	//
	// Returns:
	//   - string: The JWT token if authentication is successful.
	//   - error: An error if authentication fails, otherwise nil.
	Login(ctx context.Context, username, password string) (string, error)

	// GetUserByID retrieves a user by their ID.
	// Returns the user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be retrieved.
	//
	// Returns:
	//   - *model.User: The user model if found.
	//   - error: An error if the retrieval fails or the user is not found.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// UpdateUserByID updates a user's display name and email by their ID.
	// Returns an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be updated.
	//   - displayName: The new display name for the user.
	//   - email: The new email address for the user.
	//
	// Returns:
	//   - error: An error if the update fails, otherwise nil.
	UpdateUserByID(ctx context.Context, id, displayName, email string) error
}

type userService struct {
	userRepo        user.Repository
	passwordHashing utils.PasswordHashing
	jwtGenerator    jwtutils.JWTGenerator
}

// NewUserService creates a new instance of the  user service.
//
// Parameters:
//   - userRepo: The user repository used for database operations.
//   - passwordHashing: The password hashing utility for securing passwords.
//   - jwtGenerator: The JWT generator for creating authentication tokens.
//
// Returns:
//   - Service: A new user service instance.
func NewUserService(userRepo user.Repository, passwordHashing utils.PasswordHashing, jwtGenerator jwtutils.JWTGenerator) Service {
	return &userService{
		userRepo:        userRepo,
		passwordHashing: passwordHashing,
		jwtGenerator:    jwtGenerator,
	}
}
