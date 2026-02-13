package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/B00m3r0302/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login command requires a username \n")
	}

	username := cmd.Args[0]
	nameExists, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if username == "unknown" {
			fmt.Println("Username unknown doesn't exist and you can't login without registering\n")
			os.Exit(1)
		}
		fmt.Printf("Username doesn't exist and you can't login without %s registering\n", nameExists.Name)
		os.Exit(1)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %w \n", err)
	}
	fmt.Printf("Successfully set user to %s\n", username)
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register command requires a username\n")
	}

	username := cmd.Args[0]
	uuidVal := uuid.New()
	nowTime := time.Now().UTC()

	sameName, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		fmt.Printf("User %s already exists!\nChoose another name\n", sameName.Name)
		os.Exit(1)
	}

	databaseUser := database.CreateUserParams{
		ID:        uuidVal,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
		Name:      username,
	}

	created, err := s.db.CreateUser(context.Background(), databaseUser)
	if err != nil {
		errorMsg := fmt.Errorf("failed to create user: %w \n", err)
		return errorMsg
	}

	err = s.cfg.SetUser(created.Name)
	if err != nil {
		errorMsg := fmt.Errorf("failed to set user: %w \n", err)
		return errorMsg
	}

	fmt.Printf("Successfully set user to %s\nID: %s\nCreated at: %s\nUpdated at: %s\nName: %s\n", created.Name, created.ID, created.CreatedAt, created.UpdatedAt, created.Name)

	return nil

}

func handlerReset(s *State, cmd Command) error {
	if cmd.Name == "reset" {
		err := s.db.Reset(context.Background())
		if err != nil {
			fmt.Printf("failed to reset database: %w \n", err)
			os.Exit(1)
		}
	}
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	if cmd.Name != "users" {
		fmt.Println("Try again with users command")
		os.Exit(1)
	}
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		fmt.Printf("failed to get all users\n %w\n", err)
		os.Exit(1)
	}

	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user)
		} else {
			fmt.Printf("%s\n", user)
		}

	}
	return nil
}

type Commands struct {
	commands map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.commands[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	err := handler(s, cmd)
	if err != nil {
		errorMsg := fmt.Errorf("error running command %s: %w", cmd.Name, err)
		return errorMsg
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.commands[name] = f
}
