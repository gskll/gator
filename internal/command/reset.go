package command

import (
	"context"
	"fmt"

	"github.com/gskll/gator/internal/state"
)

func handlerReset(s *state.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting users: %w", err)
	}
	fmt.Println("Database reset.")
	return nil
}
