package main

import (
	"fmt"
	"log"
	"os"

	"github.com/B00m3r0302/aggreGATOR/internal/config"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	state := &State{
		cfg: cfg,
	}

	commands := &Commands{}

	commands.register("login", handlerLogin)

	arguments := os.Args

	if len(arguments) < 2 {
		fmt.Println("Not enough arguments. Usage: aggreGATOR <command>")
		os.Exit(1)
	}
	var commandArguments *Command

	if len(arguments) >= 3 {
		commandArguments = &Command{
			Name: arguments[1],
			Args: arguments[2:],
		}
	} else {
		commandArguments = &Command{
			Name: arguments[1],
		}
	}

	// Set the current user
	err = cfg.SetUser("kali")
	if err != nil {
		log.Fatalf("Error setting user: %v", err)
	}

	// Read the config file again
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Print the contents
	fmt.Printf("Config contents:\n")
	fmt.Printf("  DB URL: %s\n", cfg.DbUrl)
	fmt.Printf("  Current User: %s\n", cfg.CurrentUserName)
}

type State struct {
	cfg *config.Config
}
