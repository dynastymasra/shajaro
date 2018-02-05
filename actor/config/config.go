package config

import "github.com/spf13/viper"

var (
	Address string
	GinMode string
)

func setDefault() {
	viper.SetDefault("ADDRESS", ":8080")
	viper.SetDefault("GIN_MODE", "debug")
}

func InitConfig() {
	setDefault()
	viper.AutomaticEnv()

	Address = viper.GetString("ADDRESS")
	GinMode = viper.GetString("GIN_MODE")
}
