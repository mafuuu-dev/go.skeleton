package jwt_test

import (
	"backend/core/config"
	pkgjwt "backend/core/pkg/jwt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSecurity_EncodeDecode_Success(t *testing.T) {
	cfg := &config.Config{JWTSecret: "super_secret_key_123"}
	sec := pkgjwt.NewSecurity(cfg)

	claims := jwt.MapClaims{
		"user_id": "12345",
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("Failed to create jwt token: %v", err)
	}

	encoded, err := sec.Encode(tokenString)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decodedClaims, err := sec.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decodedClaims["user_id"] != "12345" {
		t.Errorf("expected user_id '12345', got '%v'", decodedClaims["user_id"])
	}
	if decodedClaims["role"] != "admin" {
		t.Errorf("expected role 'admin', got '%v'", decodedClaims["role"])
	}
}
