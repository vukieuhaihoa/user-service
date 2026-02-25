package user

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Login authenticates a user with the provided username and password.
// If authentication is successful, it generates and returns a JWT token.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - username: The username of the user attempting to log in.
//   - password: The password of the user attempting to log in.
//
// Returns:
//   - string: The JWT token if authentication is successful.
//   - error: An error if authentication fails, otherwise nil.
func (u *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	ok := u.passwordHashing.CompareHashAndPassword(user.Password, password)
	if !ok {
		return "", ErrInvalidCredentials
	}

	jwtContent := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(TokenExpirationDuration).Unix(),
	}

	token, err := u.jwtGenerator.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	return token, nil
}
