# Gator - RSS feed scraper in Go - Guided project via boot.dev
CLI RSS scraper with lightweight user management and permanent storage in Postgres database

## Requirements
- Go
- Postgres

## Installation
1. Clone repo
2. Run `go install` from the project's root
3. Add your go installation folder to PATH using `$ export PATH=$PATH:/Users/oliver.strom/go/bin/gator_boot` (on Mac/Linux)
    1 To find your go install folder, run `go list -f '{{.Target}}'` from the project's root

## Setup
1. Create a .gatorconfig.json file in your home directory

The file should contain the following format 
```json
        {
                    "db_url": "postgres://**user.name**@localhost:5432/gator?sslmode=disable",
                    "current_user_name": "",
        }
```

## Usage
Commands follow the same structure: gator <command> [arguments...]

*Example commands*
- gator register <username> : Registers a new user and saves it to the user database
- gator login <username> : Logs in the specified user. Can only be used for previously registered users
- gator users : Shows all registered users and who is currently logged in
- gator addfeed <name> <URL> : Adds the specified feed to database and registers the current user as feed follower
- gator follow <url> : Subscribes current user to the specified feed
- gator agg <interval> : Scrapes the available feeds, cycling through feeds with the specified interval. Intervals should be written as 5s, 5min, 1h etc.
