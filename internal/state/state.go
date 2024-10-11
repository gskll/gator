package state

import (
	"github.com/gskll/gator/internal/config"
	"github.com/gskll/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

func NewState(cfg config.Config, db *database.Queries) *State {
	return &State{Cfg: &cfg, Db: db}
}
