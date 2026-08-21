package dl

import (
	"context"
	"time"

	pkgdl "github.com/selyusize/my-home/pkg/dl"
	pkgdocker "github.com/selyusize/my-home/pkg/docker"
)

const serviceTimeout = 5 * time.Minute

type DLService struct {
	mgr    *pkgdl.Manager
	docker pkgdocker.Client
}

func NewDLService(engine pkgdocker.Client) (*DLService, error) {
	mgr, err := pkgdl.New()
	if err != nil {
		return nil, err
	}
	return &DLService{mgr: mgr, docker: engine}, nil
}

func (s *DLService) Status(ctx context.Context) (*pkgdl.Status, error) {
	st, err := s.mgr.Status(ctx)
	if err != nil {
		return nil, err
	}
	s.fillDocker(ctx, st)
	return st, nil
}

func (s *DLService) fillDocker(ctx context.Context, st *pkgdl.Status) {
	if s.docker == nil {
		return
	}
	if _, err := s.docker.Ping(ctx); err != nil {
		return
	}
	st.IsDockerOK = true

	if info, err := s.docker.Info(ctx); err == nil && info != nil {
		st.DockerVersion = info.ServerVersion
		st.DockerOS = info.OperatingSystem
	}

	containers, err := s.docker.ListContainers(ctx, pkgdocker.ContainerListOptions{All: true})
	if err != nil {
		return
	}
	hints := make([]pkgdl.ContainerHint, 0, len(containers))
	for _, item := range containers {
		hints = append(hints, pkgdl.ContainerHint{
			Name:   item.Name,
			Names:  item.Names,
			Image:  item.Image,
			State:  item.State,
			Labels: item.Labels,
		})
	}
	st.Services = pkgdl.DetectServices(hints)
	st.IsServiceUp = pkgdl.IsTraefikRunning(st.Services)
}

func (s *DLService) Install(ctx context.Context) error {
	return s.mgr.Install(ctx)
}

func (s *DLService) ServiceUp(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, serviceTimeout)
	defer cancel()
	return s.mgr.ServiceUp(ctx)
}

func (s *DLService) ServiceDown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, serviceTimeout)
	defer cancel()
	return s.mgr.ServiceDown(ctx)
}

func (s *DLService) Exec(ctx context.Context, workdir string, args []string) (*pkgdl.Result, error) {
	return s.mgr.Exec(ctx, workdir, args)
}

func (s *DLService) Uninstall() error {
	return s.mgr.Uninstall()
}
