package command

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gskll/gator/internal/database"
	"github.com/gskll/gator/internal/state"
)

func handlerUsers(s *state.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting users: %w", err)
	}

	for _, user := range users {
		name := user.Name
		if name == s.Cfg.CurrentUserName {
			name += " (current)"
		}
		fmt.Printf("* %s\n", name)
	}
	return nil
}

func handlerReset(s *state.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting users: %w", err)
	}
	fmt.Println("Users reset.")
	return nil
}

func handlerLogin(s *state.State, cmd Command) error {
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

func handlerRegister(s *state.State, cmd Command) error {
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

	fmt.Println("User registered:")
	printUser(user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf("* ID:	  %v\n", user.ID)
	fmt.Printf("* Name:   %v\n", user.Name)
}
