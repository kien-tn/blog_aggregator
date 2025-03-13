package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kien-tn/blog_aggregator/internal/config"
	"github.com/kien-tn/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
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

func handlerRegister(s *state, cmd command) error {
	// do something
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a username is required")
	}
	u, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	err = handlerLogin(s, command{name: "login", arguments: []string{u.Name}})
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Fprintf(os.Stdout, "User %v successfully created\n", u)
	return nil
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

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("a username is required")
	}
	_, err := s.db.GetUserByName(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("user is not found in database: %v", err)
	}
	err = s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User set to", cmd.arguments[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}
	fmt.Println("Users successfully deleted")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	fmt.Println("Users:")
	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			fmt.Printf("* %v (current)\n", u.Name)
		} else {
			fmt.Println("*", u.Name)
		}
	}
	return nil
}

func main() {

	s := &state{}
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.config = &cfg
	s.db = dbQueries
	cmd := command{name: os.Args[1], arguments: os.Args[2:]}
	cmds := commands{handlers: make(map[string]func(s *state, cmd command) error)}
	// Register the handlers
	cmds.register("update-db-url", handlerUpdateDBUrl)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerFetchFeed)
	cmds.register("addfeed", handlerAddFeed)
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
