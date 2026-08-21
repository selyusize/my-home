package docker

import "time"

// Config holds connection settings shared by docker providers.
type Config struct {
	Host       string
	APIVersion string
	CertPath   string
}

// Ping is a lightweight daemon health check.
type Ping struct {
	APIVersion     string `json:"apiVersion"`
	OSType         string `json:"osType"`
	IsExperimental bool   `json:"experimental"`
}

// Info is a simplified docker daemon summary.
type Info struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ServerVersion     string `json:"serverVersion"`
	OperatingSystem   string `json:"operatingSystem"`
	OSVersion         string `json:"osVersion"`
	OSType            string `json:"osType"`
	Architecture      string `json:"architecture"`
	NCPU              int    `json:"ncpu"`
	MemTotal          int64  `json:"memTotal"`
	Driver            string `json:"driver"`
	RootDir           string `json:"dockerRootDir"`
	Containers        int    `json:"containers"`
	ContainersRunning int    `json:"containersRunning"`
	ContainersPaused  int    `json:"containersPaused"`
	ContainersStopped int    `json:"containersStopped"`
	Images            int    `json:"images"`
}

// Port is a published container port.
type Port struct {
	IP          string `json:"ip"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort"`
	Type        string `json:"type"`
}

// Container is a simplified container list item.
type Container struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Names   []string          `json:"names"`
	Image   string            `json:"image"`
	ImageID string            `json:"imageId"`
	Command string            `json:"command"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Ports   []Port            `json:"ports"`
	Labels  map[string]string `json:"labels"`
	Created time.Time         `json:"created"`
}

// ContainerDetails is a simplified inspect result.
type ContainerDetails struct {
	Container
	Path         string    `json:"path"`
	Args         []string  `json:"args"`
	Platform     string    `json:"platform"`
	RestartCount int       `json:"restartCount"`
	PID          int       `json:"pid"`
	ExitCode     int       `json:"exitCode"`
	StartedAt    time.Time `json:"startedAt"`
	FinishedAt   time.Time `json:"finishedAt"`
}

// ContainerListOptions controls which containers are returned.
type ContainerListOptions struct {
	All   bool `json:"all"`
	Limit int  `json:"limit"`
}

// ContainerRemoveOptions controls container deletion.
type ContainerRemoveOptions struct {
	Force         bool `json:"force"`
	RemoveVolumes bool `json:"removeVolumes"`
}

// StopOptions controls graceful stop/restart timeout.
type StopOptions struct {
	Timeout *int `json:"timeout"`
}

// Image is a simplified image list item.
type Image struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repoTags"`
	RepoDigests []string          `json:"repoDigests"`
	Labels      map[string]string `json:"labels"`
	Size        int64             `json:"size"`
	Containers  int64             `json:"containers"`
	Created     time.Time         `json:"created"`
}

// ImageDetails is a simplified image inspect result.
type ImageDetails struct {
	Image
	Author       string `json:"author"`
	Comment      string `json:"comment"`
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

// Volume is a simplified named volume.
type Volume struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Scope      string            `json:"scope"`
	Labels     map[string]string `json:"labels"`
	CreatedAt  time.Time         `json:"createdAt"`
}

// Network is a simplified docker network.
type Network struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Driver     string    `json:"driver"`
	Scope      string    `json:"scope"`
	IsInternal bool      `json:"internal"`
	Created    time.Time `json:"created"`
}
