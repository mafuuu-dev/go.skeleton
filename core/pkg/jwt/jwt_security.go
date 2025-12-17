package jwt

import (
	"backend/core/config"
	"backend/core/pkg/errorsx"
	"crypto/sha256"

	"github.com/golang-jwt/jwt/v5"
	"github.com/square/go-jose/v3"
)

type Security struct {
	serverSecretKey []byte
}

func NewSecurity(cfg *config.Config) *Security {
	return &Security{
		serverSecretKey: []byte(cfg.JWTSecret),
	}
}

func (security *Security) Encode(payload string) (string, error) {
	if len(payload) == 0 {
		return "", errorsx.Errorf("Cannot encrypt empty string")
	}

	encoderKey := security.serverSecretKey
	if len(encoderKey) != 32 {
		sum := sha256.Sum256(encoderKey)
		encoderKey = sum[:]
	}

	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: encoderKey},
		nil,
	)
	if err != nil {
		return "", errorsx.Errorf("Failed to create encrypter: %v", err)
	}

	obj, err := enc.Encrypt([]byte(payload))
	if err != nil {
		return "", errorsx.Errorf("Failed to encrypt token: %v", err)
	}

	return obj.CompactSerialize()
}

func (security *Security) Decode(token string) (jwt.MapClaims, error) {
	if len(token) == 0 {
		return nil, errorsx.Errorf("Empty encoded token")
	}

	obj, err := jose.ParseEncrypted(token)
	if err != nil {
		return nil, errorsx.Errorf("Failed to parse encrypted token: %v", err)
	}

	decoderKey := security.serverSecretKey
	if len(decoderKey) != 32 {
		sum := sha256.Sum256(decoderKey)
		decoderKey = sum[:]
	}

	decryptedToken, err := obj.Decrypt(decoderKey)
	if err != nil {
		return nil, errorsx.Errorf("Failed to decrypt token: %v", err)
	}

	tokenParsed, _, err := new(jwt.Parser).ParseUnverified(string(decryptedToken), jwt.MapClaims{})
	if err != nil {
		return nil, errorsx.Errorf("Failed to parse decrypted JWT: %v", err)
	}

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorsx.Errorf("Failed to extract claims from decrypted JWT")
	}

	return claims, nil
}
