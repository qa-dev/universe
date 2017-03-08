package config

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var cwd_arg = flag.String("workdir", "", "set workdir")

func LoadConfig() *viper.Viper {
	cfg := viper.New()
	cfg.SetConfigName("config")
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

func SetTestDitectory() {
	flag.Parse()
	if *cwd_arg != "" {
		if err := os.Chdir(*cwd_arg); err != nil {
			panic(err)
		}
	}
}
