package model

import (
	"time"

	"github.com/google/uuid"
)

type TokenModel struct {
	JTI       uuid.UUID `db:"jti"`
	Token     string    `db:"token"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
