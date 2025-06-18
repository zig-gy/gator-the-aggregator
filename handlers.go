package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zig-gy/gator-the-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no username passed for login")
	}

	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("Username %s set successfuly.\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no username passed for register")
	}

	username := cmd.arguments[0]
	createdUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: username,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error writing user to config file, try logging in: %v", err)
	}

	fmt.Println("User registered succesfully")
	fmt.Println(createdUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUser(context.Background()); err != nil {
		return fmt.Errorf("error resetting user table: %v", err)
	}
	fmt.Println("Users deleted succesfully")
	return nil
}
