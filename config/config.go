package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func LoadConfig() *viper.Viper {
	cfg := viper.New()
	cfg.SetConfigFile("config.yaml")
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath(".")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	cfg.WatchConfig()
	cfg.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed: ", e.Name)
	})
	return cfg
}
