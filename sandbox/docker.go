package sandbox

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var cli client.APIClient
var clientError error

func init() {
	cli, clientError = client.NewClientWithOpts(client.WithVersion("1.39"))
	if clientError != nil {
		panic(clientError)
	}
}

func CreateContainer(path string) string {

	volumes := map[string]struct{}{path + ":/home/sandbox:rw": struct{}{}}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:     "runner",
		Cmd:       []string{"/bin/bash"},
		Tty:       true,
		OpenStdin: true,
		Volumes:   volumes,
		User:      "sandbox",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/home/sandbox",
			},
		},
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP: "0.0.0.0",
				},
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	if err = StartShell(resp.ID); err != nil {
		panic(err)
	}

	return resp.ID
}

func StartShell(id string) error {
	config := &types.ExecConfig{
		User:       "sandbox",
		Detach:     true,
		WorkingDir: "/home/sandbox",
		Cmd:        []string{"/box/shell_server.sh"},
	}

	resp, err := cli.ContainerExecCreate(context.Background(), id, *config)
	if err != nil {
		return err
	}

	startConfig := &types.ExecStartCheck{
		Detach: true,
	}

	err = cli.ContainerExecStart(context.Background(), resp.ID, *startConfig)

	return err
}

func ExecContainer(id string, lang Lang) string {
	config := &types.ExecConfig{
		Detach:     true,
		WorkingDir: "/home/sandbox",
		Cmd:        []string{"./program_runner.sh", "-c " + lang.Compiler, "-f " + lang.FileName},
	}
	resp, err := cli.ContainerExecCreate(context.Background(), id, *config)
	if err != nil {
		panic(err)
	}

	startConfig := &types.ExecStartCheck{}

	err = cli.ContainerExecStart(context.Background(), resp.ID, *startConfig)
	if err != nil {
		panic(err)
	}

	return resp.ID
}
