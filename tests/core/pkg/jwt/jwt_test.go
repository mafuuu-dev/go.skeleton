package jwt_test

import (
	"backend/core/config"
	pkgjwt "backend/core/pkg/jwt"
	"backend/internal/domain/currency/enum"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTForServer_Positive(t *testing.T) {
	j := pkgjwt.NewJWT(&config.Config{
		JWTSecret:        "server_secret_key",
		CentrifugoSecret: "socket_secret_key",
	})

	tokenStr, err := j.GenerateJWTForServer(1001, 2002, string(currency_enum.CodeUSD))
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("server_secret_key"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.EqualValues(t, 1001, claims["user_id"])
	assert.EqualValues(t, 2002, claims["account_id"])
	assert.Greater(t, int64(claims["exp"].(float64)), time.Now().Unix())
}

func TestGenerateJWTForSocket_Negative_InvalidSecret(t *testing.T) {
	j := pkgjwt.NewJWT(&config.Config{
		JWTSecret:        "real_server_secret",
		CentrifugoSecret: "real_socket_secret",
	})

	tokenStr, err := j.GenerateJWTForSocket(999)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong_socket_secret"), nil
	})

	assert.Error(t, err)
	assert.NotNil(t, parsed)
	assert.False(t, parsed.Valid)
}
