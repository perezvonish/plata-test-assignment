package quote

import (
	"context"
	"errors"
	"time"

	"perezvonish/plata-test-assignment/internal/domain/quote"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) quote.Repository {
	return &RepositoryImpl{
		pool: pool,
	}
}

func (r *RepositoryImpl) GetById(ctx context.Context, id uuid.UUID) (*quote.Quote, error) {
	query := `
       SELECT id, from_currency, to_currency, price_e8_rate, updated_at 
       FROM quotes WHERE id = $1`

	var m Model
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.Id, &m.FromCurrency, &m.ToCurrency, &m.PriceE8Rate, &m.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return m.MapToDomain(), nil
}

func (r *RepositoryImpl) GetByPair(ctx context.Context, params quote.GetByPairParams) (*quote.Quote, error) {
	query := `
       SELECT id, from_currency, to_currency, price_e8_rate, updated_at 
       FROM quotes WHERE from_currency = $1 AND to_currency = $2`

	var m Model
	err := r.pool.QueryRow(ctx, query, params.FromCurrency, params.ToCurrency).Scan(
		&m.Id, &m.FromCurrency, &m.ToCurrency, &m.PriceE8Rate, &m.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return m.MapToDomain(), nil
}

func (r *RepositoryImpl) UpdatePrice(ctx context.Context, params quote.UpdatePriceParams) error {
	query := `
       UPDATE quotes 
       SET price_e8_rate = $1, updated_at = $2 
       WHERE id = $3`

	_, err := r.pool.Exec(ctx, query, params.Price, time.Now(), params.Id)
	return err
}

func (r *RepositoryImpl) Save(ctx context.Context, q *quote.Quote) error {
	m := MapToModel(q)

	query := `
       INSERT INTO quotes (id, from_currency, to_currency, price_e8_rate, updated_at)
       VALUES ($1, $2, $3, $4, $5)
       ON CONFLICT (from_currency, to_currency) DO UPDATE SET
          price_e8_rate = EXCLUDED.price_e8_rate,
          updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		m.Id,
		m.FromCurrency,
		m.ToCurrency,
		m.PriceE8Rate,
		m.UpdatedAt,
	)
	return err
}
