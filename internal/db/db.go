package db

import (
	"github.com/selyusize/my-home/internal/config"
	pkgdb "github.com/selyusize/my-home/pkg/db"
	"github.com/selyusize/my-home/pkg/db/postgres"
)

var factory = pkgdb.NewFactory(
	postgres.Register(),
)

func New() (pkgdb.Client, error) {
	return factory.New(pkgdb.ProviderPostgres, config.Database())
}
