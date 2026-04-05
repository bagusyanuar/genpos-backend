package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	AppVersion string `mapstructure:"APP_VERSION"`
	AppEnv     string `mapstructure:"APP_ENV"`
	AppPort    string `mapstructure:"APP_PORT"`
	AppDebug   bool   `mapstructure:"APP_DEBUG"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unable to unmarshal config: %v", err)
	}

	return &conf
}
