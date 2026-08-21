package dl_test

import (
	"testing"

	"github.com/selyusize/my-home/pkg/dl"
)

func TestDetectServices(t *testing.T) {
	t.Parallel()

	hints := []dl.ContainerHint{
		{
			Name:   "dl-traefik-1",
			Image:  "traefik:v2.11",
			State:  "running",
			Labels: map[string]string{"com.docker.compose.service": "traefik"},
		},
		{
			Name:  "dl-portainer-1",
			Image: "portainer/portainer-ce:latest",
			State: "exited",
		},
		{
			Name:  "unrelated",
			Image: "nginx:latest",
			State: "running",
		},
	}

	got := dl.DetectServices(hints)
	if len(got) != 3 {
		t.Fatalf("len=%d", len(got))
	}
	if got[0].Name != "traefik" || !got[0].IsRunning || !got[0].IsPresent {
		t.Fatalf("traefik=%+v", got[0])
	}
	if got[1].Name != "portainer" || got[1].IsRunning || !got[1].IsPresent {
		t.Fatalf("portainer=%+v", got[1])
	}
	if got[2].Name != "mail" || got[2].IsRunning || got[2].IsPresent {
		t.Fatalf("mail=%+v", got[2])
	}
	if !dl.IsTraefikRunning(got) {
		t.Fatal("expected traefik running")
	}
}

func TestDetectServicesMailhog(t *testing.T) {
	t.Parallel()

	got := dl.DetectServices([]dl.ContainerHint{{
		Name:  "mailhog",
		Image: "mailhog/mailhog:latest",
		State: "running",
	}})
	if !got[2].IsRunning {
		t.Fatalf("mail=%+v", got[2])
	}
}
