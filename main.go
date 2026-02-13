package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/B00m3r0302/aggreGATOR/internal/config"
	"github.com/B00m3r0302/aggreGATOR/internal/database"
)

func main() {

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state := &State{
		cfg: cfg,
		db:  dbQueries,
	}

	cmds := &Commands{commands: make(map[string]func(*State, Command) error)}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)

	arguments := os.Args

	if len(arguments) < 2 {
		fmt.Println("Not enough arguments. Usage: aggreGATOR <command>")
		os.Exit(1)
	}
	var commandArguments Command

	if arguments[1] == "login" && len(arguments) < 3 {
		errorMessage := fmt.Errorf("Not enough arguments for command %s. Usage: aggreGATOR <command> <command arguments>", arguments[1])
		fmt.Println(errorMessage)
		os.Exit(1)
	}

	commandArguments = Command{
		Name: arguments[1],
		Args: arguments[2:],
	}

	err = cmds.Run(state, commandArguments)
	if err != nil {
		errorMessage := fmt.Errorf("Error running command: %v", err)
		fmt.Println(errorMessage)
		os.Exit(1)
	}
}

type State struct {
	db  *database.Queries
	cfg *config.Config
}
