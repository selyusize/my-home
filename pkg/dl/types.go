package dl

// Status is the local runtime snapshot for the managed dl binary.
type Status struct {
	Path              string         `json:"path"`
	IsInstalled       bool           `json:"installed"`
	Version           string         `json:"version"`
	Latest            string         `json:"latest"`
	IsUpdateAvailable bool           `json:"updateAvailable"`
	IsDockerOK        bool           `json:"dockerOk"`
	DockerVersion     string         `json:"dockerVersion"`
	DockerOS          string         `json:"dockerOs"`
	IsServiceUp       bool           `json:"serviceUp"`
	Services          []ServiceState `json:"services"`
}

// ServiceState is one dl infrastructure container (traefik, portainer, mail).
type ServiceState struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"running"`
	IsPresent bool   `json:"present"`
}

// ContainerHint is the docker snapshot used to detect dl service containers.
type ContainerHint struct {
	Name   string
	Names  []string
	Image  string
	State  string
	Labels map[string]string
}

// Result is the outcome of one dl process.
type Result struct {
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}
