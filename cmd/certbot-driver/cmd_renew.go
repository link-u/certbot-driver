package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"go.uber.org/zap"
)

func renewCerts(config *Config) {
	envs := []string{
		fmt.Sprintf("AWS_CONFIG_FILE=%s/%s", IAMMountPoint, filepath.Base(config.Aws.IamPath)),
	}
	args := []string{
		"renew",
		"--non-interactive",
		"-vvv",
		"--agree-tos",
		"--email",
		config.EmailAddress,
		"--preferred-challenges",
		"dns-01",
		"--dns-route53",
		"--dns-route53-propagation-seconds",
		"30",
	}

	if config.DryRun {
		args = append(args, "--dry-run", "--test-cert")
	}

	err := createRunner(context.Background(), config, args, envs).Run()
	if err != nil {
		log.Fatal("Failed to run container", zap.Error(err))
	}
}
