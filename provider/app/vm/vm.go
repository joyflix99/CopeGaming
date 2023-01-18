package vm

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func loadenv(rel_path string) {
	path, _ := os.Getwd()
	path += rel_path
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func StartVM(id string, appName string, videoRelayPort, audioRelayPort, syncPort int) error {
	log.Printf("[%s] Spinning off VM\n", id)

	params := []string{
		id,
		strconv.Itoa(videoRelayPort),
		strconv.Itoa(audioRelayPort),
		strconv.Itoa(syncPort),
		appName,
	}

	// Load Appconf .env file
	loadenv("/appconf/" + appName + ".env")
	env := os.Environ()
	env = append(env, "videoport="+strconv.Itoa(videoRelayPort))
	env = append(env, "audioport="+strconv.Itoa(audioRelayPort))
	env = append(env, "wsport="+strconv.Itoa(syncPort))
	//Print params
	for _, value := range params {
		log.Printf("[%s] params\n", value)

	}

	// Start Docker Container
	println("Start Docker Container")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	//Pass Enviroment
	gpuOpts := opts.GpuOpts{}
	gpuOpts.Set("all")
	mountpath, _ := os.Getwd()
	mountpath += "/../appvm/apps/" + appName
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "provider_appvm:latest",
		Tty:   true,
		Env:   env,
	}, &container.HostConfig{
		Resources: container.Resources{DeviceRequests: gpuOpts.Value()},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountpath,
				Target: "/appvm/app",
			},
		},
	}, nil, nil, id)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID, nil
}

func StopVM(id, appName string, container_id string) error {
	log.Printf("[%s] Stopping VM\n", id)

	params := []string{
		id,
		appName,
	}
	cmd := exec.Command("./stopVM.sh", params...)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
