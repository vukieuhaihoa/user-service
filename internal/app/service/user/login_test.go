package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockJWT "github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils/mocks"
	mockPasswordHashing "github.com/vukieuhaihoa/bookmark-libs/pkg/utils/mocks"
	"github.com/vukieuhaihoa/user-service/internal/app/model"
	mockUserRepo "github.com/vukieuhaihoa/user-service/internal/app/repository/user/mocks"
)

var ErrCannotGenerateToken = errors.New("cannot generate token")

func TestService_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepo     func(ctx context.Context) *mockUserRepo.Repository
		setupMockPasswordHash func(t *testing.T) *mockPasswordHashing.PasswordHashing
		setupMockJWTGen       func(t *testing.T) *mockJWT.JWTGenerator

		inputUsername string
		inputPassword string

		expectedError error

		expectedOutput string
	}{
		{
			name: "Login successfully",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					Base: model.Base{
						ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
					},
					Username: "testuser",
					Password: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", // hash for "password123"
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("CompareHashAndPassword", "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", "password123").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				jwtMock.On("GenerateToken", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					if claims["sub"] != "de305d54-75b4-431b-adb2-eb6b9e546099" {
						return false
					}

					if _, ok := claims["iat"].(int64); !ok {
						return false
					}

					if _, ok := claims["exp"].(int64); !ok {
						return false
					}

					return true
				})).Return("mocked_jwt_token", nil)
				return jwtMock
			},

			inputUsername: "testuser",
			inputPassword: "password123",

			expectedOutput: "mocked_jwt_token",
		},
		{
			name: "Fail to get user by username",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "nonexistentuser").Return(nil, ErrInvalidCredentials)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				return mockPasswordHashing.NewPasswordHashing(t)
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			inputUsername: "nonexistentuser",
			inputPassword: "somepassword",

			expectedError: ErrInvalidCredentials,
		},
		{
			name: "Invalid password",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					Base: model.Base{
						ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
					},
					Username: "testuser",
					Password: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", // hash for "password123"
				}, nil)
				return repoMock
			},
			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("CompareHashAndPassword", "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", "wrongpassword").Return(false)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			inputUsername: "testuser",
			inputPassword: "wrongpassword",

			expectedError: ErrInvalidCredentials,
		},
		{
			name: "Fail to generate JWT token",

			setupMockUserRepo: func(ctx context.Context) *mockUserRepo.Repository {
				repoMock := mockUserRepo.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					Base: model.Base{
						ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
					},
					Username: "testuser",
					Password: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", // hash for "password123"
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mockPasswordHashing.PasswordHashing {
				hashingMock := mockPasswordHashing.NewPasswordHashing(t)
				hashingMock.On("CompareHashAndPassword", "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", "password123").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				jwtMock.On("GenerateToken", mock.Anything).Return("", ErrCannotGenerateToken)
				return jwtMock
			},

			inputUsername: "testuser",
			inputPassword: "password123",

			expectedError: ErrCannotGenerateToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			userRepoMock := tc.setupMockUserRepo(ctx)
			passwordHashingMock := tc.setupMockPasswordHash(t)
			jwtGenMock := tc.setupMockJWTGen(t)

			userService := NewUserService(userRepoMock, passwordHashingMock, jwtGenMock)

			res, err := userService.Login(ctx, tc.inputUsername, tc.inputPassword)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)

		})
	}
}
