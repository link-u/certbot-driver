package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/link-u/certbot-bot/internal/bot"
	"go.uber.org/zap"
)

func renewCerts(config *Config) {
	runner := bot.Runner{
		Context: context.Background(),
		ContainerConfig: container.Config{
			Image: ImageName,
			Cmd:   []string{""},
		},
	}
	err := runner.Run()
	if err != nil {
		log.Fatal("Failed to run container", zap.Error(err))
	}
}
