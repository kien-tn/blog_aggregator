package main

import (
	"fmt"
	"os"

	"github.com/kien-tn/blog_aggregator/internal/config"
)

type state struct {
	config *config.Config
}
type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(s *state, cmd command) error
}

func (c *commands) register(name string, handler func(s *state, cmd command) error) {
	c.handlers[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func handlerUpdateDBUrl(s *state, cmd command) error {
	// do something
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a db url is required")
	}
	err := s.config.SetDBUrl(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error setting db url: %w", err)
	}
	fmt.Println("DB URL successfully updated")
	return nil
}

func handlerLogin(s *state, cmd command) error {
	// do something
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a username is required")
	}
	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User set to", cmd.arguments[0])
	return nil
}

func main() {
	s := &state{}
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s.config = &cfg
	cmd := command{name: os.Args[1], arguments: os.Args[2:]}
	cmds := commands{handlers: make(map[string]func(s *state, cmd command) error)}
	// Register the handlers
	cmds.register("update-db-url", handlerUpdateDBUrl)
	cmds.register("login", handlerLogin)
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing argument")
		os.Exit(1)
	}
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config.ReadCfgFile()
}
