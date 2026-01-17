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
		return "", errorsx.New("Cannot encrypt empty string")
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
		return "", errorsx.Wrap(err, "Failed to create encrypter")
	}

	obj, err := enc.Encrypt([]byte(payload))
	if err != nil {
		return "", errorsx.Wrap(err, "Failed to encrypt token")
	}

	return obj.CompactSerialize()
}

func (security *Security) Decode(token string) (jwt.MapClaims, error) {
	if len(token) == 0 {
		return nil, errorsx.New("Empty encoded token")
	}

	obj, err := jose.ParseEncrypted(token)
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to parse encrypted token")
	}

	decoderKey := security.serverSecretKey
	if len(decoderKey) != 32 {
		sum := sha256.Sum256(decoderKey)
		decoderKey = sum[:]
	}

	decryptedToken, err := obj.Decrypt(decoderKey)
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to decrypt token")
	}

	tokenParsed, _, err := new(jwt.Parser).ParseUnverified(string(decryptedToken), jwt.MapClaims{})
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to parse decrypted JWT")
	}

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorsx.New("Failed to extract claims from decrypted JWT")
	}

	return claims, nil
}
