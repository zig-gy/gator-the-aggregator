package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no username passed for login")
	}

	username := cmd.arguments[0]
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Username %s set successfuly.\n", username)
	return nil
}