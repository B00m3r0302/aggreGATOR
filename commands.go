package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
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
	if len(cmd.Args) != 1 {
		fmt.Println("Try again with a time argument after agg command")
		os.Exit(1)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	// Run immediately
	err = ScrapeFeeds(s, context.Background())
	if err != nil {
		return fmt.Errorf("failed to scrape feeds: %w", err)
	}

	// Then run every tick
	for range ticker.C {
		err = ScrapeFeeds(s, context.Background())
		if err != nil {
			return fmt.Errorf("failed to scrape feeds: %w", err)
		}
	}

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

func handlerBrowse(s *State, cmd Command) error {
	limit := int32(2)
	if len(cmd.Args) > 0 {
		var parsedLimit int
		_, err := fmt.Sscanf(cmd.Args[0], "%d", &parsedLimit)
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(parsedLimit)
	}

	userId, err := s.db.GetUserId(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user id: %w", err)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: userId,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found")
		return nil
	}

	fmt.Printf("\nShowing %d posts:\n", len(posts))
	for _, post := range posts {
		fmt.Printf("\n---\n")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.Format("2006-01-02 15:04:05"))
		}
	}
	fmt.Printf("\n---\n")

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

func parsePublishedAt(pubDate string) sql.NullTime {
	// Try common RSS date formats
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"2006-01-02T15:04:05Z07:00", // ISO 8601
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}

	// If we can't parse the date, return null
	return sql.NullTime{Valid: false}
}

func ScrapeFeeds(s *State, ctx context.Context) error {
	fetchFeedList, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	if len(fetchFeedList) == 0 {
		fmt.Println("No feeds to fetch")
		return nil
	}

	for _, feedData := range fetchFeedList {
		fmt.Printf("\nFetching feed: %s\n", feedData.Name)
		fmt.Printf("URL: %s\n", feedData.Url)

		feed, err := rss.FetchFeed(ctx, feedData.Url)
		if err != nil {
			fmt.Printf("Error fetching feed: %v\n", err)
			continue
		}

		fmt.Printf("Found %d posts\n", len(feed.Channel.Item))

		savedCount := 0
		for _, item := range feed.Channel.Item {
			publishedAt := parsePublishedAt(item.PubDate)

			postParams := database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Url:         item.Link,
				Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
				PublishedAt: publishedAt,
				FeedID:      feedData.ID,
			}

			_, err := s.db.CreatePost(ctx, postParams)
			if err != nil {
				// Check if it's a duplicate URL error
				if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
					// Silently ignore duplicate posts
					continue
				}
				// Log other errors
				fmt.Printf("  Error saving post '%s': %v\n", item.Title, err)
			} else {
				savedCount++
			}
		}

		fmt.Printf("Saved %d new posts\n", savedCount)

		markFeedParams := database.MarkLastFeedFetchedParams{
			LastFetchedAt: time.Now().UTC(),
			ID:            feedData.ID,
		}

		result, err := s.db.MarkLastFeedFetched(ctx, markFeedParams)
		if err != nil {
			fmt.Printf("Error marking feed as fetched: %v\n", err)
		} else {
			fmt.Printf("Marked feed as fetched at: %v\n", result.LastFetchedAt)
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
