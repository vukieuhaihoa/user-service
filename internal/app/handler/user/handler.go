// Package user provides HTTP handlers for user-related operations.
// It includes handlers for creating users, logging in, retrieving user profiles,
// and updating user profiles using the Gin web framework.
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/user-service/internal/app/service/user"
)

// Handler defines the interface for user-related HTTP handlers.
// It provides methods for handling user creation requests using the Gin web framework.
type Handler interface {
	// CreateUser is a Gin framework handler that creates a new user.
	// It processes HTTP requests and returns the created user or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	CreateUser(c *gin.Context)

	// Login is a Gin framework handler that authenticates a user and returns a JWT token.
	// It processes HTTP requests and returns the token or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	Login(c *gin.Context)

	// GetProfile is a Gin framework handler that retrieves the profile of the authenticated user.
	// It processes HTTP requests and returns the user profile or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	GetProfile(c *gin.Context)

	// UpdateProfile is a Gin framework handler that updates the profile of the authenticated user.
	// It processes HTTP requests and returns a success message or an error.
	//
	// Parameters:
	//   - c: The Gin context containing the HTTP request and response
	UpdateProfile(c *gin.Context)
}

// userHandler is the concrete implementation of the Handler interface.
type userHandler struct {
	userSvc user.Service
}

// NewUser creates a new instance of the User handler.
// It accepts a user service implementation and returns a handler
// that can process HTTP requests for user creation and login.
//
// Parameters:
//   - userSvc: The user service used for user-related operations
//
// Returns:
//   - Handler: A new user handler instance
func NewUserHandler(userSvc user.Service) Handler {
	return &userHandler{userSvc: userSvc}
}
