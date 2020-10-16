package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/link-u/certbot-bot/internal/bot"
	"go.uber.org/zap"
)

func createCerts(config *Config, domains []string) {
	log := zap.L()
	certDir, err := filepath.Abs(config.Cert.Directory)
	if err != nil {
		log.Fatal("Failed to calculate abs path", zap.String("path", config.Cert.Directory), zap.Error(err))
	}
	certDir = filepath.Clean(certDir)
	err = os.MkdirAll(certDir, 0777)
	if err != nil {
		log.Fatal("Failed to create cert dir", zap.String("path", config.Cert.Directory), zap.Error(err))
	}

	envs := []string{
		fmt.Sprintf("AWS_CONFIG_FILE=%s", config.Aws.IamPath),
	}
	args := []string{
		"run",
		"--service-ports",
		"--rm certbot",
		"certonly",
		"-vvv",
		"--agree-tos",
		fmt.Sprintf("--email %s", config.EmailAddress),
		"--keep",
		"--preferred-challenges dns-01",
		"--non-interactive",
		"--dns-route53",
		"--dns-route53-propagation-seconds 30",
	}
	runner := bot.Runner{
		Context: context.Background(),
		ContainerConfig: container.Config{
			Image: ImageName,
			Cmd:   args,
			Env:   envs,
		},
		HostConfig: container.HostConfig{
			AutoRemove:  true,
			NetworkMode: "host",
			Mounts: []mount.Mount{
				{
					Type:     "bind",
					Source:   certDir,
					Target:   "/etc/letsencrypt",
					ReadOnly: false,
				},
			},
		},
	}
	err = runner.Run()
	if err != nil {
		log.Fatal("Failed to run container", zap.Error(err))
	}
}
