package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/gskll/gator/internal/command"
	"github.com/gskll/gator/internal/config"
	"github.com/gskll/gator/internal/database"
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
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	dbQueries := database.New(db)
	state := state.NewState(cfg, dbQueries)

	cmds := command.NewCommands()
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)

	cmd := command.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
