package config

import "github.com/spf13/viper"

var (
	Address string
	GinMode string

	DatabaseUsername string
	DatabasePassword string
	DatabaseAddress  string
	DatabaseName     string
	DatabaseLog      bool
	DatabaseMaxOpen  int
	DatabaseMaxIdle  int
)

func setDefault() {
	viper.SetDefault("ADDRESS", ":8080")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("DATABASE_ADDRESS", "localhost:5432")
	viper.SetDefault("DATABASE_USERNAME", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "root")
	viper.SetDefault("DATABASE_NAME", "actor")
	viper.SetDefault("DATABASE_LOG", true)
	viper.SetDefault("DATABASE_MAX_OPEN", 25)
	viper.SetDefault("DATABASE_MAX_IDLE", 2)
}

func InitConfig() {
	setDefault()
	viper.AutomaticEnv()

	Address = viper.GetString("ADDRESS")
	GinMode = viper.GetString("GIN_MODE")

	DatabaseAddress = viper.GetString("DATABASE_ADDRESS")
	DatabaseUsername = viper.GetString("DATABASE_USERNAME")
	DatabasePassword = viper.GetString("DATABASE_PASSWORD")
	DatabaseName = viper.GetString("DATABASE_NAME")
	DatabaseLog = viper.GetBool("DATABASE_LOG")
	DatabaseMaxOpen = viper.GetInt("DATABASE_MAX_OPEN")
	DatabaseMaxIdle = viper.GetInt("DATABASE_MAX_IDLE")
}
