package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/link-u/certbot-bot/internal/bot"
	"go.uber.org/zap"
)

func createRunner(ctx context.Context, config *Config, args []string, envs []string) *bot.Runner {
	log := zap.L()
	certDir, err := filepath.Abs(config.Cert.Directory)
	if err != nil {
		log.Fatal("Failed to calculate abs path", zap.String("path", config.Cert.Directory), zap.Error(err))
	}
	certDir = filepath.Clean(certDir)
	iamPath, err := filepath.Abs(config.Aws.IamPath)
	if err != nil {
		log.Fatal("Failed to calculate abs path", zap.String("path", config.Cert.Directory), zap.Error(err))
	}
	iamPath = filepath.Clean(iamPath)
	iamDir := filepath.Dir(iamPath)
	err = os.MkdirAll(certDir, 0777)
	if err != nil {
		log.Fatal("Failed to create cert dir", zap.String("path", config.Cert.Directory), zap.Error(err))
	}

	runner := bot.Runner{
		Context: ctx,
		ContainerConfig: container.Config{
			Image: ImageName,
			Cmd:   args,
			Env:   envs,
		},
		HostConfig: container.HostConfig{
			AutoRemove:  false,
			NetworkMode: "host",
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					Source:   certDir,
					Target:   CertsMountPoint,
					ReadOnly: false,
					BindOptions: &mount.BindOptions{
						Propagation: mount.PropagationShared,
					},
				},
				{
					Type:     mount.TypeBind,
					Source:   iamDir,
					Target:   IAMMountPoint,
					ReadOnly: false,
					BindOptions: &mount.BindOptions{
						Propagation: mount.PropagationRShared,
					},
				},
			},
		},
	}
	return &runner
}
