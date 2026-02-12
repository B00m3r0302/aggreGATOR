package main

import "fmt"

type Command struct {
	Name string
	Args []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login command requires a username or at least one value after the 'login' command")
	}

	username := cmd.Args[0]
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("Successfully set user to %s\n", username)
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
