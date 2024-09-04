package postgres

type Config struct {
	DBURL          string `env:"DB_URL,notEmpty"`
	MigrationsPath string `env:"MIGRATIONS_PATH,notEmpty"`
}

func NewConfig(dburl string, migrPath string) Config {
	return Config{DBURL: dburl, MigrationsPath: migrPath}
}
