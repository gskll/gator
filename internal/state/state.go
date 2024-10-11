package state

import "github.com/gskll/gator/internal/config"

type State struct {
	Cfg *config.Config
}

func NewState(cfg config.Config) *State {
	return &State{Cfg: &cfg}
}
