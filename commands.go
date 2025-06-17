package main

import (
	"fmt"

	"github.com/zig-gy/gator-the-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if err := c.cmds[cmd.name](s, cmd); err != nil {
		return fmt.Errorf("error running command: %v", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no username passed for login")
	}

	username := cmd.arguments[0]
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Username %s set successfuly.", username)
	return nil
}

