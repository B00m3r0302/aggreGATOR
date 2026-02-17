package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/B00m3r0302/aggreGATOR/internal/database"
	"github.com/B00m3r0302/aggreGATOR/internal/rss"
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
	nameExists, err := s.db.GetUserId(context.Background(), username)
	if err != nil {
		if username == "unknown" {
			fmt.Println("Username unknown doesn't exist and you can't login without registering\n")
			os.Exit(1)
		}
		fmt.Printf("Username doesn't exist and you can't login without %s registering\n", nameExists)
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

	sameName, err := s.db.GetUserId(context.Background(), username)
	if err == nil {
		fmt.Printf("User %s already exists!\nChoose another name\n", sameName)
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

func handlerAgg(s *State, cmd Command) error {
	// if len(cmd.Args) == 0 {
	//	fmt.Println("Try again with an argument after agg command")
	//	os.Exit(1)
	// }

	// rss.FetchFeed(context.Background(), cmd.Args[0])
	rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	return nil
}

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	// Check args length
	if len(cmd.Args) != 2 {
		fmt.Errorf("addfeed command requires two arguments: feed name and feed url")
	}

	// Get current user_id
	currentUser := s.cfg.CurrentUserName
	userId, err := s.db.GetUserId(context.Background(), currentUser)
	if err != nil {
		fmt.Printf("failed to get id for %s\n %w\n", currentUser, err)
		os.Exit(1)
	}

	// add feed struct
	addFeedStruct := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    userId,
	}

	// insert record
	da_feed, err := s.db.AddFeed(context.Background(), addFeedStruct)
	if err != nil {
		fmt.Printf("failed to add feed\n %w\n", err)
		os.Exit(1)
	}

	err = handlerFollow(s, Command{"follow", []string{cmd.Args[1]}})
	if err != nil {
		fmt.Printf("failed to follow feed\n %w\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added feed %v\n", da_feed)
	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.db.GetUserFeeds(context.Background())
	if err != nil {
		fmt.Printf("failed to get feeds\n %w\n", err)
		os.Exit(1)
	}
	for _, feed := range feeds {
		fmt.Printf("%v\n", feed)
	}
	return nil
}

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		fmt.Errorf("follow command requires two arguments: feed url")
	}

	userId, err := s.db.GetUserId(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Printf("failed to get id for %s\n %w\n", s.cfg.CurrentUserName, err)
		os.Exit(1)
	}

	url := cmd.Args[0]
	feed, err := s.db.QueryByUrl(context.Background(), url)
	if err != nil {
		fmt.Printf("failed to get feed id for %s\n %w\n", url, err)
		os.Exit(1)
	}

	feedId := feed.ID

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userId,
		FeedID:    feedId,
	}

	feedback, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Printf("failed to create feed follow\n %w\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully followed feed %v for %v\n", feedback.FeedName, feedback.UserName)

	return nil
}

func handlerFollowing(s *State, cmd Command) error {
	userId, err := s.db.GetUserId(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Printf("failed to get id for %s\n %w\n", s.cfg.CurrentUserName, err)
		os.Exit(1)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), userId)
	if err != nil {
		fmt.Printf("failed to get feed follows for %s\n %w\n", s.cfg.CurrentUserName, err)
		os.Exit(1)
	}

	for _, f := range follows {
		fmt.Println(f.FeedName)
	}
	return nil
}

func handlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		fmt.Errorf("unfollow command requires one argument: feed url")
	}

	userId, err := s.db.GetUserId(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Printf("failed to get id for %s\n %w\n", s.cfg.CurrentUserName, err)
		os.Exit(1)
	}

	url := cmd.Args[0]
	urlData, err := s.db.QueryByUrl(context.Background(), url)
	if err != nil {
		fmt.Printf("failed to get feed id for %s\n %w\n", url, err)
		os.Exit(1)
	}
	feedId := urlData.ID

	params := database.UnfollowFeedParams{
		FeedID: feedId,
		UserID: userId,
	}

	err = s.db.UnfollowFeed(context.Background(), params)
	return nil
}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		username, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		return handler(s, cmd, username)
	}
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
