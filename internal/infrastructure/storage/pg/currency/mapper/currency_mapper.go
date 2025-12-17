package pg_currency_mapper

import (
	"backend/core/pkg/pgscan"
	"backend/internal/domain/currency/entity"
	"backend/internal/infrastructure/storage/model"
)

func ToEntity(row pgscan.Scannable) (currency_entity.Currency, error) {
	var currency currency_entity.Currency
	err := row.Scan(
		&currency.ID,
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.Precision,
		&currency.IsCrypto,
		&currency.IsActive,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	return currency, err
}

func ToModel(currency currency_entity.Currency) model.CurrencyModel {
	return model.CurrencyModel{
		ID:        currency.ID,
		Code:      string(currency.Code),
		Name:      currency.Name,
		Symbol:    currency.Symbol,
		Precision: currency.Precision,
		IsCrypto:  currency.IsCrypto,
		IsActive:  currency.IsActive,
		CreatedAt: currency.CreatedAt,
		UpdatedAt: currency.UpdatedAt,
	}
}
