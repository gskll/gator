package command

import (
	"errors"

	"github.com/gskll/gator/internal/state"
)

var ErrHandlerNotFound = errors.New("handler not found")

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	registeredCommands map[string]func(*state.State, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		registeredCommands: make(map[string]func(*state.State, Command) error),
	}
}

func (c *Commands) RegisterCommands() {
	c.register("reset", handlerReset)

	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("users", handlerUsers)

	c.register("agg", handlerAgg)
	c.register("feeds", handlerFeeds)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("browse", middlewareLoggedIn(handlerBrowse))
}

func (c *Commands) register(name string, f func(*state.State, Command) error) error {
	c.registeredCommands[name] = f
	return nil
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	handler, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return ErrHandlerNotFound
	}
	return handler(s, cmd)
}
