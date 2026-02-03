package config

type Config struct {
	Server            ServerConfig
	Postgres          PostgresConfig
	JobWorker         JobWorkerConfig
	ExchangeApiConfig ExchangeApiConfig
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
	Port int64  `env:"SERVER_PORT" envDefault:"8080"`
}

type JobWorkerConfig struct {
	WorkerCount int `env:"JOB_WORKER_COUNT" envDefault:"1"`
}

type PostgresConfig struct {
	Host              string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port              int    `env:"POSTGRES_PORT" envDefault:"5432"`
	Username          string `env:"POSTGRES_USERNAME" envDefault:"postgres"`
	Password          string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	DatabaseName      string `env:"POSTGRES_DATABASE_NAME" envDefault:"postgres"`
	ConnectRetryCount int    `env:"POSTGRES_DATABASE_CONNECT_RETRY_COUNT" envDefault:"10"`
}

type ExchangeApiConfig struct {
	URL string `env:"EXCHANGE_API_URL" required:"true"`
	Key string `env:"EXCHANGE_API_KEY" required:"true"`
}
