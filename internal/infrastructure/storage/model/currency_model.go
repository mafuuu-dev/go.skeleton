package model

import (
	"time"
)

type CurrencyModel struct {
	ID        int64     `db:"id"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	Symbol    string    `db:"symbol"`
	Precision int32     `db:"precision"`
	IsCrypto  bool      `db:"is_crypto"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
