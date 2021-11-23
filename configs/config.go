package configs

import (
	"github.com/spf13/viper"
)
const CFGPATH = "CFG_PATH"

func getDefaultPath()string{
	viper.AutomaticEnv()
	return viper.GetString(CFGPATH)
}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(getDefaultPath())
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
