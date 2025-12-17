package currency_repository

import (
	"backend/internal/domain/currency/entity"
	"backend/internal/domain/currency/enum"
)

type Repository interface {
	GetCurrencyByCode(code currency_enum.Code) (*currency_entity.Currency, error)
}
