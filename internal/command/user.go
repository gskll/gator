package command

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gskll/gator/internal/database"
	"github.com/gskll/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	username := cmd.Args[0]

	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User doesn't exist: %w", err)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}

	fmt.Printf("Logged in: %s\n", user.Name)
	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	username := cmd.Args[0]

	now := time.Now()
	user, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{ID: uuid.New(), Name: username, UpdatedAt: now, CreatedAt: now},
	)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}

	fmt.Printf("Registered: %+v", user)
	return nil
}
