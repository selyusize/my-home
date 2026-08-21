package git

import (
	pkggit "github.com/selyusize/my-home/pkg/git"
	"github.com/selyusize/my-home/pkg/git/github"
	"github.com/selyusize/my-home/pkg/git/gitlab"
)

var factory = pkggit.NewFactory(
	github.Register(),
	gitlab.Register(),
)
