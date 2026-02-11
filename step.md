# Gator Multi-User CLI - Implementation Steps

## Overview
This document describes the implementation of the Gator multi-user CLI application with JSON-based configuration management.

## Step 1: Create the Config File

Created `~/.gatorconfig.json` in the home directory with initial structure:

```json
{
  "db_url": "postgres://example"
}
```

Note: The `current_user_name` field is intentionally omitted and will be set by the application.

## Step 2: Initialize Go Module

The Go module was already initialized:

```bash
go mod init github.com/B00m3r0302/aggreGATOR
```

This created `go.mod` with Go version 1.24.9.

## Step 3: Create Main Function

Created `main.go` with a basic main function that imports the config package.

## Step 4: Create Internal Config Package

Created the directory structure:
- `internal/config/`

Implemented `internal/config/config.go` with the following components:

### Constants
```go
const configFileName = ".gatorconfig.json"
```

### Exported Types
```go
type Config struct {
    DbUrl           string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}
```

The struct includes JSON tags for proper serialization/deserialization.

### Exported Functions

#### Read() Function
```go
func Read() (*Config, error)
```
- Reads the config file from `~/.gatorconfig.json`
- Uses `getConfigFilePath()` to get the full path
- Opens and decodes the JSON file into a Config struct
- Returns pointer to Config and any errors encountered

#### SetUser() Method
```go
func (c *Config) SetUser(username string) error
```
- Updates the `CurrentUserName` field
- Writes the updated config back to the file
- Returns any errors from the write operation

### Internal Helper Functions

#### getConfigFilePath()
```go
func getConfigFilePath() (string, error)
```
- Uses `os.UserHomeDir()` to get the home directory
- Appends the config filename constant
- Returns the full path or an error

#### write()
```go
func write(cfg Config) error
```
- Gets the config file path
- Marshals the Config struct to JSON
- Writes the JSON to the file with 0644 permissions
- Returns any errors encountered

## Step 5: Update Main Function

Updated `main.go` to demonstrate the config functionality:

1. **Read the config file**
   ```go
   cfg, err := config.Read()
   ```

2. **Set the current user to "kali"**
   ```go
   err = cfg.SetUser("kali")
   ```

3. **Read the config file again**
   ```go
   cfg, err = config.Read()
   ```

4. **Print the config contents**
   ```go
   fmt.Printf("Config contents:\n")
   fmt.Printf("  DB URL: %s\n", cfg.DbUrl)
   fmt.Printf("  Current User: %s\n", cfg.CurrentUserName)
   ```

## Step 6: Test the Implementation

Ran the program:
```bash
go run main.go
```

### Output
```
Config contents:
  DB URL: postgres://example
  Current User: kali
```

### Verified Config File
The `~/.gatorconfig.json` file was successfully updated:
```json
{"db_url":"postgres://example","current_user_name":"kali"}
```

## Key Implementation Details

### Error Handling
- All file operations include proper error handling
- Errors are wrapped with context using `fmt.Errorf` with `%w` verb
- Main function uses `log.Fatalf` for fatal errors

### File Operations
- Config file is located in the user's home directory (`~/.gatorconfig.json`)
- Read operations use `json.NewDecoder` for streaming JSON parsing
- Write operations use `json.Marshal` followed by `os.WriteFile`
- File permissions set to 0644 (readable by all, writable by owner)

### Package Organization
- Config logic is encapsulated in the `internal/config` package
- Internal package prevents external imports (Go best practice)
- Helper functions are unexported (lowercase) for internal use only
- Public API consists of: `Config` struct, `Read()` function, and `SetUser()` method

## Architecture Benefits

1. **Single Source of Truth**: One JSON file manages both database connection and current user
2. **Multiplayer Support**: Multiple users can use the same database on a single device
3. **Persistence**: Configuration persists between program runs
4. **Encapsulation**: Config logic is isolated in its own package
5. **Error Handling**: Robust error handling with descriptive messages

## Next Steps

The application can now be extended with:
- User management commands (register, login, etc.)
- Database connection logic using the `db_url`
- Additional CLI commands for multi-user functionality
- PostgreSQL integration for data persistence