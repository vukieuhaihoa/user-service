package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	service "github.com/vukieuhaihoa/user-service/internal/app/service/user"
)

type loginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser001"`
	Password string `json:"password" binding:"required,gte=8" example:"my_SECURE_password123@"`
}

// Login generates a Gin framework handler that authenticates a user and returns a JWT token.
// @Summary      User login
// @Description  Authenticate a user and return a JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        credentials  body      loginRequest  true  "User credentials"
// @Success      200          {object}  object{data=string,message=string}
// @Failure      400          {object}  object{message=string}
// @Failure      401          {object}  object{message=string}
// @Failure      500          {object}  object{message=string}
// @Router       /v1/users/login [post]
func (u *userHandler) Login(c *gin.Context) {
	// Implementation for user login handler goes here
	input := &loginRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, common.InputFieldError(err))
		return
	}

	token, err := u.userSvc.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusBadRequest, common.Message{
			Message: err.Error(),
		})
		return
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusBadRequest, common.Message{
			Message: "invalid username or password",
		})
		return
	case errors.Is(err, nil):
	default:
		log.Error().
			Str("operation", "Login").
			Err(err).
			Msg("service return error when login user")
		c.JSON(http.StatusInternalServerError, common.InternalErrorResponse)
		return

	}

	c.JSON(http.StatusOK, &common.SuccessResponse[string]{
		Data:    token,
		Message: "Logged in successfully!",
	})
}
