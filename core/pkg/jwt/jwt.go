package jwt

import (
	"backend/core/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	serverSecretKey []byte
	socketSecretKet []byte
}

func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		serverSecretKey: []byte(cfg.JWTSecret),
		socketSecretKet: []byte(cfg.CentrifugoSecret),
	}
}

func (j *JWT) GenerateJWTForServer(userID int64, accountID int64, currencyCode string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       userID,
		"account_id":    accountID,
		"currency_code": currencyCode,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.serverSecretKey)
}

func (j *JWT) GenerateJWTForSocket(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.socketSecretKet)
}
