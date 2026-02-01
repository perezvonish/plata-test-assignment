-- Таблица для хранения курсов валют
CREATE TABLE IF NOT EXISTS quotes (
    pair        VARCHAR(10) PRIMARY KEY,
    from_curr   VARCHAR(3) NOT NULL,
    to_curr     VARCHAR(3) NOT NULL,
    rate        NUMERIC(18, 8) NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для отслеживания задач на обновление
CREATE TABLE IF NOT EXISTS update_jobs (
    id              UUID PRIMARY KEY,
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    status          VARCHAR(20) NOT NULL, -- pending, completed, failed
    from_curr       VARCHAR(3) NOT NULL,
    to_curr         VARCHAR(3) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);