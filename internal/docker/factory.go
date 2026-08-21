package docker

import (
	pkgdocker "github.com/selyusize/my-home/pkg/docker"
	"github.com/selyusize/my-home/pkg/docker/engine"
)

var factory = pkgdocker.NewFactory(
	engine.Register(),
)
