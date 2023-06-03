package main

import (
	"github.com/spf13/viper"
	"log"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("[loadConfig] load config err: %v", err)
	}
	Run(conf.Tunnels)
}

type Config struct {
	Tunnels []*Tunnel `json:"tunnels"`
}

func loadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("conf")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var conf Config
	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
