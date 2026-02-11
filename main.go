package main

import (
	"fmt"
	"log"

	"github.com/B00m3r0302/aggreGATOR/internal/config"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
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
