package quote

import (
	"context"
	"errors"
	"fmt"
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

func (r *RepositoryImpl) Save(ctx context.Context, params quote.SaveParams) (*quote.Quote, error) {
	id := uuid.New()
	now := time.Now().UTC()

	query := `
       INSERT INTO quotes (id, from_currency, to_currency, price_e8_rate, created_at, updated_at)
       VALUES ($1, $2, $3, $4, $5, $6)
       ON CONFLICT (from_currency, to_currency) DO UPDATE SET 
          price_e8_rate = EXCLUDED.price_e8_rate,
          updated_at = EXCLUDED.updated_at
       RETURNING id, from_currency, to_currency, price_e8_rate, created_at, updated_at`

	var q quote.Quote
	err := r.pool.QueryRow(ctx, query,
		id,
		params.FromCurrency,
		params.ToCurrency,
		0, // Начальный курс
		now,
		now,
	).Scan(
		&q.Id,
		&q.FromCurrency,
		&q.ToCurrency,
		&q.PriceE8Rate,
		&q.CreatedAt,
		&q.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save and scan quote: %w", err)
	}

	return &q, nil
}
