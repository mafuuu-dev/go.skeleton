package token_entity

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	JTI       uuid.UUID `json:"jti"`
	Token     string    `json:"token"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
