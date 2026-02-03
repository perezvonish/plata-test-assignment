CREATE TABLE IF NOT EXISTS quotes (
    id             UUID PRIMARY KEY,
    from_currency  VARCHAR(3) NOT NULL,
    to_currency    VARCHAR(3) NOT NULL,
    price_e8_rate  BIGINT NOT NULL DEFAULT 0,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_quotes_currency_pair UNIQUE (from_currency, to_currency)
);

CREATE TABLE IF NOT EXISTS update_jobs (
    id              UUID PRIMARY KEY,
    quote_id        UUID NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    status          VARCHAR(20) NOT NULL,
    retry_count     BIGINT NOT NULL DEFAULT 0,
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON update_jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_quote_id ON update_jobs(quote_id);