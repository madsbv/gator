package state

import (
	"github.com/madsbv/gator/internal/config"
	"github.com/madsbv/gator/internal/database"
)

type State struct {
	Config config.Config
	Db     database.Queries
}
