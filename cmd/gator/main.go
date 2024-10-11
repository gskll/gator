package main

import (
	"log"
	"os"

	"github.com/gskll/gator/internal/command"
	"github.com/gskll/gator/internal/config"
	"github.com/gskll/gator/internal/state"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	state := state.NewState(cfg)

	cmds := command.NewCommands()
	cmds.Register("login", command.HandlerLogin)

	cmd := command.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
