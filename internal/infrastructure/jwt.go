package infrastructure

import (
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils"
)

// CreateJWTProviders initializes and returns JWT generator and validator.
// Returns:
//   - jwtutils.JWTGenerator: The initialized JWT generator
//   - jwtutils.JWTValidator: The initialized JWT validator
func CreateJWTProviders() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGenerator, err := jwtutils.NewJWTGenerator("./private_key.pem")
	common.HandlerError(err)

	jwtValidator, err := jwtutils.NewJWTValidator("./public_key.pem")
	common.HandlerError(err)

	return jwtGenerator, jwtValidator
}
