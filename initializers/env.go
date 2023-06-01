package initializers

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Port       string `env:"PORT" envDefault:"8080"`
	LogFile    string `env:"LOG_FILE" envDefault:"gin.log"`
	DBSync     bool   `env:"DB_SYNC" envDefault:"false"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"postgres"`
}

func LoadEnv() (Env, error) {
	var envCfg = &Env{}

	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		return Env{}, err
	}

	if err := env.Parse(envCfg); err != nil {
		return Env{}, err
	}
	return *envCfg, nil
}
