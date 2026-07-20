package main

import (
	"fmt"

	"github.com/liebeSonne/shortlink/internal/config"
)

func initConfig() (config.Config, error) {
	cfg, err := config.LoadConfig(appID, envPrefix)
	if err != nil {
		return config.Config{}, fmt.Errorf("error get config: %s", err.Error())
	}
	return cfg, nil
}
