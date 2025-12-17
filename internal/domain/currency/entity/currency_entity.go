package currency_entity

import (
	"backend/internal/domain/currency/enum"
	"time"
)

type Currency struct {
	ID        int64              `json:"id"`
	Code      currency_enum.Code `json:"code"`
	Name      string             `json:"name"`
	Symbol    string             `json:"symbol"`
	Precision int32              `json:"precision"`
	IsCrypto  bool               `json:"is_crypto"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
