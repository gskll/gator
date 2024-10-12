package command

import (
	"context"
	"fmt"

	"github.com/gskll/gator/internal/database"
	"github.com/gskll/gator/internal/state"
)

func middlewareLoggedIn(
	handler func(*state.State, Command, database.User) error,
) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error getting user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
