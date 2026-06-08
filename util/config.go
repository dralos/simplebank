package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `env:"DBDRIVER" envDefault:"postgres"`
	DBSource            string        `env:"DBSOURCE" envDefault:"postgresql://app_user:secret@localhost:5432/simple_bank?sslmode=disable"`
	ServerAddress       string        `env:"SERVERADDRESS" envDefault:"0.0.0.0:8080"`
	TokenSymmetricKey   string        `env:"TOKENSYMMETRICKEY" envDefault:"a_random_symmetric_key_123456789"`
	AccessTokenDuration time.Duration `env:"ACCESSTOKENDURATION" envDefault:"15m"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
