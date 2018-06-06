package config

import "github.com/spf13/viper"

var (
	Address          string
	GinMode          string
	RetryDuration    int
	MaxRetryInterval int

	DatabaseUsername string
	DatabasePassword string
	DatabaseAddress  string
	DatabaseName     string
	DatabaseLog      bool
	DatabaseMaxOpen  int
	DatabaseMaxIdle  int

	KongURL string
)

func setDefault() {
	viper.SetDefault("ADDRESS", ":8080")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("RETRY_DURATION", 30)
	viper.SetDefault("MAX_RETRY_INTERVAL", 5)

	viper.SetDefault("DATABASE_ADDRESS", "localhost:5432")
	viper.SetDefault("DATABASE_USERNAME", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "root")
	viper.SetDefault("DATABASE_NAME", "actor")
	viper.SetDefault("DATABASE_LOG", true)
	viper.SetDefault("DATABASE_MAX_OPEN", 25)
	viper.SetDefault("DATABASE_MAX_IDLE", 2)

	viper.SetDefault("KONG_URL", "http://localhost:8001")
}

func InitConfig() {
	setDefault()
	viper.AutomaticEnv()

	Address = viper.GetString("ADDRESS")
	GinMode = viper.GetString("GIN_MODE")
	RetryDuration = viper.GetInt("RETRY_DURATION")
	MaxRetryInterval = viper.GetInt("MAX_RETRY_INTERVAL")

	DatabaseAddress = viper.GetString("DATABASE_ADDRESS")
	DatabaseUsername = viper.GetString("DATABASE_USERNAME")
	DatabasePassword = viper.GetString("DATABASE_PASSWORD")
	DatabaseName = viper.GetString("DATABASE_NAME")
	DatabaseLog = viper.GetBool("DATABASE_LOG")
	DatabaseMaxOpen = viper.GetInt("DATABASE_MAX_OPEN")
	DatabaseMaxIdle = viper.GetInt("DATABASE_MAX_IDLE")

	KongURL = viper.GetString("KONG_URL")
}
