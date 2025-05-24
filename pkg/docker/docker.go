package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerData struct {
	Containers []ContainerInfo      `json:"containers"`
	Images     []ImageInfo          `json:"images"`
	Volumes    []VolumeInfo         `json:"volumes"`
	Networks   []NetworkInfo        `json:"networks"`
	Info       map[string]string    `json:"info"`
	Stats      map[string]StatsInfo `json:"stats"`
}

type ContainerInfo struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Image      string            `json:"image"`
	ImageID    string            `json:"imageId"`
	Command    string            `json:"command"`
	Created    time.Time         `json:"created"`
	State      string            `json:"state"`
	Status     string            `json:"status"`
	Ports      []types.Port      `json:"ports"`
	Labels     map[string]string `json:"labels"`
	HostConfig HostConfigInfo    `json:"hostConfig"`
	Mounts     []MountInfo       `json:"mounts"`
	Networks   []string          `json:"networks"`
}

type StatsInfo struct {
	CPUPercentage    float64 `json:"cpuPercentage"`
	MemoryPercentage float64 `json:"memoryPercentage"`
	MemoryUsage      int64   `json:"memoryUsage"`
	MemoryLimit      int64   `json:"memoryLimit"`
	NetworkRx        int64   `json:"networkRx"`
	NetworkTx        int64   `json:"networkTx"`
}

type HostConfigInfo struct {
	NetworkMode   string `json:"networkMode"`
	Privileged    bool   `json:"privileged"`
	RestartPolicy string `json:"restartPolicy"`
}

type MountInfo struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

type ImageInfo struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repoTags"`
	RepoDigests []string          `json:"repoDigests"`
	Created     time.Time         `json:"created"`
	Size        int64             `json:"size"`
	Labels      map[string]string `json:"labels"`
}

type VolumeInfo struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
	CreatedAt  string            `json:"createdAt"`
	Status     map[string]string `json:"status"`
}

type NetworkInfo struct {
	ID         string                          `json:"id"`
	Name       string                          `json:"name"`
	Driver     string                          `json:"driver"`
	Scope      string                          `json:"scope"`
	IPAM       IPAMInfo                        `json:"ipam"`
	Internal   bool                            `json:"internal"`
	Attachable bool                            `json:"attachable"`
	Labels     map[string]string               `json:"labels"`
	Containers map[string]ContainerNetworkInfo `json:"containers"`
}

type IPAMInfo struct {
	Driver  string            `json:"driver"`
	Options map[string]string `json:"options"`
	Config  []IPAMConfig      `json:"config"`
}

type IPAMConfig struct {
	Subnet  string `json:"subnet"`
	Gateway string `json:"gateway"`
}

type ContainerNetworkInfo struct {
	Name        string `json:"name"`
	EndpointID  string `json:"endpointId"`
	MacAddress  string `json:"macAddress"`
	IPv4Address string `json:"ipv4Address"`
	IPv6Address string `json:"ipv6Address"`
}

func CollectDockerData(ctx context.Context, dockerHost string) (interface{}, error) {
	var cli *client.Client
	var err error

	if dockerHost != "" {
		os.Setenv("DOCKER_HOST", dockerHost)
		cli, err = client.NewClientWithOpts(client.WithHost(dockerHost), client.WithAPIVersionNegotiation())
	} else {
		cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	result := DockerData{}

	info, err := cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker info: %w", err)
	}

	result.Info = map[string]string{
		"ID":                info.ID,
		"Name":              info.Name,
		"OS":                info.OperatingSystem,
		"OSType":            info.OSType,
		"Architecture":      info.Architecture,
		"KernelVersion":     info.KernelVersion,
		"DockerVersion":     info.ServerVersion,
		"Containers":        fmt.Sprintf("%d", info.Containers),
		"ContainersRunning": fmt.Sprintf("%d", info.ContainersRunning),
		"ContainersPaused":  fmt.Sprintf("%d", info.ContainersPaused),
		"ContainersStopped": fmt.Sprintf("%d", info.ContainersStopped),
		"Images":            fmt.Sprintf("%d", info.Images),
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result.Stats = make(map[string]StatsInfo)
	for _, c := range containers {
		inspect, err := cli.ContainerInspect(ctx, c.ID)
		if err != nil {
			fmt.Printf("Warning: failed to inspect container %s: %v\n", c.ID, err)
			continue
		}

		var mounts []MountInfo
		for _, m := range inspect.Mounts {
			mounts = append(mounts, MountInfo{
				Type:        string(m.Type),
				Source:      m.Source,
				Destination: m.Destination,
				Mode:        m.Mode,
				RW:          m.RW,
			})
		}

		networkNames := []string{}
		for networkName := range inspect.NetworkSettings.Networks {
			networkNames = append(networkNames, networkName)
		}

		restartPolicy := "no"
		if inspect.HostConfig.RestartPolicy.Name != "" {
			restartPolicy = string(inspect.HostConfig.RestartPolicy.Name)
		}

		hostConfig := HostConfigInfo{
			NetworkMode:   string(inspect.HostConfig.NetworkMode),
			Privileged:    inspect.HostConfig.Privileged,
			RestartPolicy: restartPolicy,
		}

		containerName := c.Names[0]
		if len(containerName) > 0 && containerName[0] == '/' {
			containerName = containerName[1:]
		}

		containerInfo := ContainerInfo{
			ID:         c.ID,
			Name:       containerName,
			Image:      c.Image,
			ImageID:    c.ImageID,
			Command:    c.Command,
			Created:    time.Unix(c.Created, 0),
			State:      c.State,
			Status:     c.Status,
			Ports:      c.Ports,
			Labels:     c.Labels,
			HostConfig: hostConfig,
			Mounts:     mounts,
			Networks:   networkNames,
		}

		result.Containers = append(result.Containers, containerInfo)

		if c.State == "running" {
			stats, err := cli.ContainerStats(ctx, c.ID, false)
			if err == nil {
				var statsResult types.StatsJSON
				err = json.NewDecoder(stats.Body).Decode(&statsResult)
				stats.Body.Close()

				if err == nil {
					cpuDelta := float64(statsResult.CPUStats.CPUUsage.TotalUsage) - float64(statsResult.PreCPUStats.CPUUsage.TotalUsage)
					systemDelta := float64(statsResult.CPUStats.SystemUsage) - float64(statsResult.PreCPUStats.SystemUsage)
					cpuPercent := 0.0
					if systemDelta > 0.0 && cpuDelta > 0.0 {
						cpuPercent = (cpuDelta / systemDelta) * float64(len(statsResult.CPUStats.CPUUsage.PercpuUsage)) * 100.0
					}

					memoryUsage := statsResult.MemoryStats.Usage
					memoryLimit := statsResult.MemoryStats.Limit
					memoryPercent := 0.0
					if memoryLimit > 0 {
						memoryPercent = float64(memoryUsage) / float64(memoryLimit) * 100.0
					}

					var networkRx, networkTx int64
					for _, network := range statsResult.Networks {
						networkRx += int64(network.RxBytes)
						networkTx += int64(network.TxBytes)
					}

					result.Stats[c.ID] = StatsInfo{
						CPUPercentage:    cpuPercent,
						MemoryPercentage: memoryPercent,
						MemoryUsage:      int64(memoryUsage),
						MemoryLimit:      int64(memoryLimit),
						NetworkRx:        networkRx,
						NetworkTx:        networkTx,
					}
				}
			}
		}
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	for _, img := range images {
		result.Images = append(result.Images, ImageInfo{
			ID:          img.ID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Created:     time.Unix(img.Created, 0),
			Size:        img.Size,
			Labels:      img.Labels,
		})
	}

	volumeListOptions := volume.ListOptions{
		Filters: filters.NewArgs(),
	}
	volumes, err := cli.VolumeList(ctx, volumeListOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	for _, vol := range volumes.Volumes {
		if vol != nil {
			status := make(map[string]string)
			if vol.Status != nil {
				for k, v := range vol.Status {
					if str, ok := v.(string); ok {
						status[k] = str
					} else {
						status[k] = fmt.Sprintf("%v", v)
					}
				}
			}

			result.Volumes = append(result.Volumes, VolumeInfo{
				Name:       vol.Name,
				Driver:     vol.Driver,
				Mountpoint: vol.Mountpoint,
				Labels:     vol.Labels,
				Scope:      vol.Scope,
				CreatedAt:  vol.CreatedAt,
				Status:     status,
			})
		}
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	for _, nw := range networks {
		network, err := cli.NetworkInspect(ctx, nw.ID, types.NetworkInspectOptions{})
		if err != nil {
			fmt.Printf("Warning: failed to inspect network %s: %v\n", nw.ID, err)
			continue
		}

		var ipamConfigs []IPAMConfig
		ipamDriver := ""
		ipamOptions := map[string]string{}

		if network.IPAM.Driver != "" {
			ipamDriver = network.IPAM.Driver
		}

		if network.IPAM.Options != nil {
			for k, v := range network.IPAM.Options {
				ipamOptions[k] = v
			}
		}

		for _, conf := range network.IPAM.Config {
			ipamConfigs = append(ipamConfigs, IPAMConfig{
				Subnet:  conf.Subnet,
				Gateway: conf.Gateway,
			})
		}

		containerMap := make(map[string]ContainerNetworkInfo)
		for id, container := range network.Containers {
			containerMap[id] = ContainerNetworkInfo{
				Name:        container.Name,
				EndpointID:  container.EndpointID,
				MacAddress:  container.MacAddress,
				IPv4Address: container.IPv4Address,
				IPv6Address: container.IPv6Address,
			}
		}

		result.Networks = append(result.Networks, NetworkInfo{
			ID:     network.ID,
			Name:   network.Name,
			Driver: network.Driver,
			Scope:  network.Scope,
			IPAM: IPAMInfo{
				Driver:  ipamDriver,
				Options: ipamOptions,
				Config:  ipamConfigs,
			},
			Internal:   network.Internal,
			Attachable: network.Attachable,
			Labels:     network.Labels,
			Containers: containerMap,
		})
	}

	return result, nil
}

func CheckCredentials(ctx context.Context, dockerHost string) (bool, error) {
	var cli *client.Client
	var err error

	if dockerHost != "" {
		cli, err = client.NewClientWithOpts(client.WithHost(dockerHost), client.WithAPIVersionNegotiation())
	} else {
		cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	}
	if err != nil {
		return false, err
	}
	defer cli.Close()

	_, err = cli.Ping(ctx)
	return err == nil, err
}

func TestConnection(ctx context.Context, host string) (string, error) {
	var opts []client.Opt

	if host != "" {
		opts = append(opts, client.WithHost(host))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return "", err
	}
	defer cli.Close()

	cli.NegotiateAPIVersion(ctx)

	serverVersion, err := cli.ServerVersion(ctx)
	if err != nil {
		return "", err
	}

	return serverVersion.Version, nil
}
