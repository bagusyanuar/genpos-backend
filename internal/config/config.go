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
 
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	JWTIssuer            string `mapstructure:"JWT_ISSUER"`
	JWTExpiration        int    `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshSecret     string `mapstructure:"JWT_SECRET_REFRESH"`
	JWTRefreshExpiration int    `mapstructure:"JWT_EXPIRATION_REFRESH"`

	LogFile       string `mapstructure:"LOG_FILE"`
	LogMaxSize    int    `mapstructure:"LOG_MAX_SIZE"`
	LogMaxBackups int    `mapstructure:"LOG_MAX_BACKUPS"`
	LogMaxAge     int    `mapstructure:"LOG_MAX_AGE"`
	LogCompress   bool   `mapstructure:"LOG_COMPRESS"`
	LogLevel      string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")

	// Logger Defaults
	viper.SetDefault("LOG_FILE", ".logs/app.log")
	viper.SetDefault("LOG_MAX_SIZE", 10)
	viper.SetDefault("LOG_MAX_BACKUPS", 3)
	viper.SetDefault("LOG_MAX_AGE", 7)
	viper.SetDefault("LOG_COMPRESS", true)
	viper.SetDefault("LOG_LEVEL", "info")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unable to unmarshal config: %v", err)
	}

	return &conf
}
