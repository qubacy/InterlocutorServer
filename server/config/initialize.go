package config

import (
	"github.com/spf13/viper"
)

func Initialize() error {
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
