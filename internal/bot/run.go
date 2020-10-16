package bot

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type Runner struct {
	Context         context.Context
	ContainerConfig container.Config
	HostConfig      container.HostConfig
	NetworkConfig   network.NetworkingConfig
}

// from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringRunes(n int) string {
	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	l := len(runes)
	for i := range b {
		b[i] = runes[rand.Intn(l)]
	}
	return string(b)
}

func (runner *Runner) Run() error {
	log := zap.L()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	{
		result, err := cli.ImagePull(runner.Context, runner.ContainerConfig.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		log.Info("Image pulled", zap.String("image-name", runner.ContainerConfig.Image))
		_, _ = io.Copy(os.Stderr, result)
		defer func() {
			_ = result.Close()
		}()
	}

	containerName := "certbot-" + RandStringRunes(64)
	resp, err := cli.ContainerCreate(runner.Context, &runner.ContainerConfig, &runner.HostConfig, &runner.NetworkConfig, containerName)
	if err != nil {
		return err
	}
	log.Info("Container created", zap.String("id", resp.ID), zap.String("name", containerName))

	if err := cli.ContainerStart(runner.Context, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil
	}
	log.Info("Container started", zap.String("id", resp.ID), zap.String("name", containerName))
	running := true
	defer func() {
		if running {
			err := cli.ContainerStop(runner.Context, resp.ID, nil)
			if err == nil {
				log.Info("Container stopped", zap.String("id", resp.ID), zap.String("name", containerName))
			} else {
				log.Warn("Failed to stop container", zap.String("id", resp.ID), zap.String("name", containerName), zap.Error(err))
				err = cli.ContainerKill(runner.Context, resp.ID, "SIGKILL")
				if err == nil {
					log.Info("Container killed", zap.String("id", resp.ID), zap.String("name", containerName))
				} else {
					log.Warn("Failed to kill container", zap.String("id", resp.ID), zap.String("name", containerName), zap.Error(err))
				}
			}
		}
		err := cli.ContainerRemove(runner.Context, resp.ID, types.ContainerRemoveOptions{})
		if err != nil {
			log.Warn("Failed to remove container", zap.String("id", resp.ID), zap.String("name", containerName), zap.Error(err))
		}
		log.Info("Container removed", zap.String("id", resp.ID), zap.String("name", containerName))
	}()

	if _, err = cli.ContainerWait(runner.Context, resp.ID); err != nil {
		return err
	}
	running = false

	out, err := cli.ContainerLogs(runner.Context, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil
	}
	var buff bytes.Buffer
	_, _ = io.Copy(os.Stdout, out)
	log.Info("Executed", zap.ByteString("log", buff.Bytes()))
	return nil
}
