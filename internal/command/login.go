package command

import (
	"fmt"

	"github.com/gskll/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	username := cmd.Args[0]

	err := s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}

	fmt.Printf("Logged in: %s\n", username)
	return nil
}
