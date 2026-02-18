# aggreGATOR

A command-line RSS feed aggregator built in Go that fetches, stores, and displays RSS feeds from multiple sources.

## Features

- 👤 **User Management** - Register and login users
- 📰 **Feed Management** - Add, follow, and unfollow RSS feeds
- 🔄 **Automatic Fetching** - Continuously scrape feeds at configurable intervals
- 💾 **Post Storage** - Save posts to PostgreSQL database with duplicate detection
- 📖 **Browse Posts** - View posts from all feeds you follow
- 🕒 **Smart Fetching** - Fetches the least recently updated feeds first

## Prerequisites

Before you can run aggreGATOR, you'll need to have the following installed on your system:

- **Go 1.21 or higher** - [Download and install Go](https://go.dev/doc/install)
- **PostgreSQL** - [Download and install PostgreSQL](https://www.postgresql.org/download/)

## Installation

### Option 1: Install using `go install` (Recommended)

This will install the `aggreGATOR` CLI directly to your `$GOPATH/bin`:

```bash
go install github.com/B00m3r0302/aggreGATOR@latest
```

Make sure your `$GOPATH/bin` is in your `PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Option 2: Build from source

1. Clone the repository:
```bash
git clone https://github.com/B00m3r0302/aggreGATOR.git
cd aggreGATOR
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o aggreGATOR
```

### Database Setup

1. Create a PostgreSQL database:
```bash
createdb gator
```

2. Install goose for migrations:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3. Run database migrations:
```bash
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
cd ../..
```

**Note:** Update the connection string if your PostgreSQL credentials are different.

## Configuration

The application stores configuration in `~/.gatorconfig.json`. This file is automatically created when you first register or login.

### Setting Up the Config File

You have two options:

**Option 1: Let the app create it (Recommended)**
Simply run the register command, and the config file will be created automatically:
```bash
aggreGATOR register <your_username>
```

**Option 2: Create it manually**
Create `~/.gatorconfig.json` with the following content:
```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

**Important:** If your PostgreSQL setup uses different credentials, update the `db_url` in the config file:
- Replace `postgres:postgres` with `username:password`
- Replace `gator` with your database name if different
- Update `localhost:5432` if your PostgreSQL runs on a different host/port

## Quick Start

Get up and running in 4 easy steps:

1. **Register a user:**
```bash
aggreGATOR register yourname
```

2. **Add an RSS feed:**
```bash
aggreGATOR addfeed "Hacker News" https://hnrss.org/newest
```

3. **Start the aggregator (fetches every minute):**
```bash
aggreGATOR agg 1m
```

4. **Browse your posts (in another terminal):**
```bash
aggreGATOR browse 10
```

That's it! The aggregator will continuously fetch new posts from your feeds.

## Usage

### User Management

**Register a new user:**
```bash
aggreGATOR register <username>
```

**Login as an existing user:**
```bash
aggreGATOR login <username>
```

**List all users:**
```bash
aggreGATOR users
```

> **Note:** If you built from source instead of using `go install`, use `aggreGATOR` instead of `aggreGATOR`

### Feed Management

**Add a new feed:**
```bash
aggreGATOR addfeed <feed_name> <feed_url>
```
Example:
```bash
aggreGATOR addfeed "Hacker News" https://hnrss.org/newest
```

**List all feeds:**
```bash
aggreGATOR feeds
```

**Follow a feed:**
```bash
aggreGATOR follow <feed_url>
```

**Unfollow a feed:**
```bash
aggreGATOR unfollow <feed_url>
```

**List feeds you're following:**
```bash
aggreGATOR following
```

### Aggregation

**Start the aggregator:**
```bash
aggreGATOR agg <time_between_requests>
```

The aggregator will continuously fetch feeds at the specified interval. Duration examples:
- `1s` - Every second
- `1m` - Every minute
- `1h` - Every hour
- `30s` - Every 30 seconds

Example:
```bash
aggreGATOR agg 1m
```

This runs in a loop and will:
1. Fetch the least recently updated feed
2. Parse all posts from the feed
3. Save new posts to the database (skipping duplicates)
4. Mark the feed as fetched
5. Wait for the specified duration
6. Repeat

Press `Ctrl+C` to stop the aggregator.

### Browsing Posts

**View recent posts:**
```bash
aggreGATOR browse [limit]
```

Examples:
```bash
aggreGATOR browse      # Show 2 posts (default)
aggreGATOR browse 10   # Show 10 most recent posts
```

Posts are displayed from all feeds you follow, ordered by publication date.

### Database Management

**Reset the database:**
```bash
aggreGATOR reset
```
⚠️ This will delete ALL data from the database.

## Project Structure

```
aggreGATOR/
├── internal/
│   ├── config/        # Configuration file handling
│   ├── database/      # Generated database queries (sqlc)
│   └── rss/          # RSS feed fetching and parsing
├── sql/
│   ├── queries/      # SQL query definitions
│   └── schema/       # Database migrations
├── commands.go       # Command handlers
├── main.go          # Application entry point
└── README.md        # This file
```

## Database Schema

### Tables

- **users** - User accounts
- **feeds** - RSS feed sources
- **feed_follows** - Many-to-many relationship between users and feeds
- **posts** - Individual posts from feeds

### Key Features

- Posts have unique URLs to prevent duplicates
- Feeds track `last_fetched_at` for efficient polling
- Cascade deletes ensure referential integrity

## Example Workflow

```bash
# 1. Register and login
aggreGATOR register alice

# 2. Add some feeds
aggreGATOR addfeed "Hacker News" https://hnrss.org/newest
aggreGATOR addfeed "Go Blog" https://go.dev/blog/feed.atom

# 3. Start aggregating (fetch every minute)
aggreGATOR agg 1m

# 4. In another terminal, browse posts
aggreGATOR browse 5
```

## RSS Feed Compatibility

The aggregator supports standard RSS 2.0 feeds and handles:
- Multiple date formats (RFC1123, RFC822, ISO 8601)
- HTML entity decoding in titles and descriptions
- Null/missing fields gracefully

## Development

### Adding New SQL Queries

1. Add your query to `sql/queries/<name>.sql`
2. Run `sqlc generate` to generate Go code
3. Use the generated functions in your commands

### Adding New Migrations

1. Create a new migration file in `sql/schema/`
2. Follow the naming convention: `N_description.sql`
3. Use goose format with `-- +goose Up` and `-- +goose Down`
4. Run migrations with goose

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

Built as a learning project following the [Boot.dev](https://boot.dev) curriculum.
