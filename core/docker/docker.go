package docker

/*
docker package is responsible for managing all interaciones with docker SDK.
It is intended that it is abstracting all the neccarrosy functions
to pull and run a container and keep track of it status.
*/

// TODO: Add image pull options: TAGS .
import (
	"bahamut/core/types"
	"context"
	"io"
	"log"
	"os"
	"time"

	dTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type (
	Runtime struct {
		ContainerID string
	}
	Config struct {
		Name          string
		AttachStdin   bool
		AttachStout   bool
		AttachStderr  bool
		Cmd           []string
		Image         string
		Memory        int64
		Disk          int64
		Env           []string
		RestartPolicy types.RestartPolicy
	}
	Docker struct {
		Client      *client.Client
		Config      Config
		ContainerId string
		Runtime     Runtime
		StartTime   time.Time
		EndTime     time.Time
	}
	DockerResult struct {
		Error       error
		Action      string
		ContainerId string
		Result      string
	}
)

func NewConig(Name string, Image string, Env []string, RestartPolicy types.RestartPolicy) *Config {
	return &Config{
		Name:          Name,
		Image:         Image,
		Env:           Env,
		RestartPolicy: RestartPolicy,
	}
}

func (c *Config) AttachCmd(Cmd []string) {
	c.Cmd = append(c.Cmd, Cmd...)
}

func (c *Config) AttachEnv(Env []string) {
	c.Env = append(c.Env, Env...)
}

func NewDocker(c *Config) *Docker {
	dCli, _ := client.NewClientWithOpts(client.FromEnv)
	docker := &Docker{
		Client: dCli,
		Config: *c,
	}
	return docker
}

//Run the container
func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(
		ctx, d.Config.Image, dTypes.ImagePullOptions{},
	)
	if err != nil {
		log.Printf("Error pulling image: %v", err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)
	rP := string(d.Config.RestartPolicy)
	restarPolicy := container.RestartPolicy{Name: rP}

	resources := container.Resources{
		Memory: d.Config.Memory,
	}

	containerConfig := container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}

	hostConfig := container.HostConfig{
		RestartPolicy:   restarPolicy,
		Resources:       resources,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(
		ctx, &containerConfig, &hostConfig, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %v\n", d.Config.Name)

		return DockerResult{Error: err}
	}

	containerStartedErr := d.Client.ContainerStart(ctx, resp.ID, dTypes.ContainerStartOptions{})

	if containerStartedErr != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, containerStartedErr)
		return DockerResult{Error: containerStartedErr}
	}

	d.Runtime.ContainerID = resp.ID

	out, err := d.Client.ContainerLogs(
		ctx, resp.ID, dTypes.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	if err != nil {
		log.Printf("Error getting logs for container %v: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerId: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (d *Docker) Stop() DockerResult {
	log.Printf("Attempting to stop container")
	ctx := context.Background()

	err := d.Client.ContainerStop(ctx, d.Runtime.ContainerID, nil)
	if err != nil {
		panic(err)
	}

	rOptions := dTypes.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}

	err = d.Client.ContainerRemove(
		ctx, d.Runtime.ContainerID, rOptions,
	)

	if err != nil {
		panic(err)
	}

	return DockerResult{
		Action: "stop", Result: "success", Error: nil,
	}

}
