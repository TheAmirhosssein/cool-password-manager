package testdocker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/necmettindev/randomstring"
)

func GenerateName() (string, error) {
	opts := randomstring.GenerationOptions{
		Length: 10,
	}

	return randomstring.GenerateString(opts)
}

func StopAndRemoveContainer(ctx context.Context, containerName, volumeName string) error {
	// Try to find the container
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	var containerID string
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				containerID = c.ID
				break
			}
		}
	}

	if containerID == "" {
		return errors.New("container not found, nothing to stop")
	}

	// Stop the container
	fmt.Println("Stopping container...")
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return err
	}

	// Remove the container
	fmt.Println("Removing container...")
	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return err
	}

	// Remove the volume
	fmt.Println("Removing volume...")
	if err := cli.VolumeRemove(ctx, volumeName, true); err != nil {
		return err
	}

	return nil
}

func StartPostgresContainer(ctx context.Context, containerName, volumeName string) (string, error) {
	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	// Pull PostgreSQL image
	imageName := "postgres:latest"
	fmt.Println("Pulling PostgreSQL image...")
	// _, err = cli.ImagePull(ctx, imageName, image.PullOptions{})
	// if err != nil {
	// 	return "", err
	// }

	// Define container configuration
	containerConfig := &container.Config{
		Image: imageName,
		Env: []string{
			"POSTGRES_USER=admin",
			"POSTGRES_PASSWORD=admin",
			"POSTGRES_DB=testdb",
		},
		ExposedPorts: nat.PortSet{
			"5432/tcp": struct{}{},
		},
	}

	// Define host configuration (port binding)
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "",
				},
			},
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Binds:         []string{volumeName + ":/var/lib/postgresql/data"},
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		return "", err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	fmt.Println("PostgreSQL container started successfully!")

	// Wait for a few seconds to allow PostgreSQL to initialize
	time.Sleep(3 * time.Second)

	// Get container details
	containerInfo, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", err
	}

	fmt.Printf("PostgreSQL is running at: %s:%s\n", containerInfo.NetworkSettings.IPAddress, "5432")

	return containerInfo.NetworkSettings.Ports["5432/tcp"][0].HostPort, nil
}

func GetTestDB(ctx context.Context, port string) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("postgres://admin:admin@localhost:%v/testdb", port)
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return pool, nil
}
