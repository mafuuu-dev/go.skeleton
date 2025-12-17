package pg_currency_query

import (
	"backend/core/pkg/pgscan"
	"backend/core/pkg/query"
	"backend/internal/domain/currency/entity"
	"backend/internal/domain/currency/enum"
	"backend/internal/infrastructure/storage/pg/currency/mapper"
)

type GetCurrencyByCode struct {
	*query.Query
	code currency_enum.Code
}

func NewGetCurrencyByCode(factory *query.Factory) *GetCurrencyByCode {
	return &GetCurrencyByCode{
		Query: factory.Instance(),
	}
}

func (q *GetCurrencyByCode) SetCode(code currency_enum.Code) *GetCurrencyByCode {
	q.code = code
	return q
}

func (q *GetCurrencyByCode) Execute() (*currency_entity.Currency, error) {
	row := q.QueryRow(q, string(q.code))

	return pgscan.ScanOne[currency_entity.Currency](row, pg_currency_mapper.ToEntity)
}

func (q *GetCurrencyByCode) Sql() string {
	return `
		SELECT 
		    id,
		    code,
		    name,
		    symbol,
		    precision,
		    is_crypto,
		    is_active,
		    created_at, 
		    updated_at
		FROM currencies
		WHERE code = $1
	`
}
