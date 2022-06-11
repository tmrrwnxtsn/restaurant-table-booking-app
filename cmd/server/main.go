package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/config"
)

var flagConfig = flag.String("config", "./configs/default.yml", "path to config file")

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatalf("failed to load config data: %s", err)
	}

	fmt.Println(cfg.DSN, cfg.BindAddr, cfg.LogLevel)
}
