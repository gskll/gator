# Gator - RSS Feed Aggregator

Gator is a command-line RSS feed aggregator written in Go. It allows users to manage feeds, follow/unfollow them, and browse posts from their followed feeds.

## Dependencies

- Go 1.23
- PostgreSQL

### External Go packages:

github.com/lib/pq
github.com/google/uuid

## Installation

- Ensure you have Go installed on your system.
- `git clone https://github.com/gskll/gator.git`
- `cd gator`
- `go mod tidy`
- `build -o gator`

## Configuration

Create a configuration file `~/.gatorconfig.json` with the following structure:
```
{
  "db_url": "postgres://username:password@localhost:5432/database_name"
}
```
Replace the database URL with your PostgreSQL connection string.

## Database Setup

Create a PostgreSQL database for the project.

Run the migrations using goose (make sure goose is installed):
- `goose postgres "your_database_url" up`

(Install `goose` directly or use `make install-tools`)

## Usage

Run the CLI with the following syntax:

`gator <command> [args...]`

Available commands:

`reset`: Reset the application state
`login <username>`: Log in as a user
`register <username>`: Register a new user
`users`: List all users
`agg <time_between_reqs>`: Aggregate feeds at specified intervals
`feeds`: List all feeds
`addfeed <name> <feed_url>`: Add a new feed (requires login)
`follow <feed_url>`: Follow a feed (requires login)
`unfollow <feed_url>`: Unfollow a feed (requires login)
`following`: List followed feeds (requires login)
`browse [number_of_posts]`: Browse posts from followed feeds (requires login)

Example usage:
./gator register johndoe
./gator login johndoe
./gator addfeed "Tech News" https://example.com/tech-rss
./gator follow https://example.com/tech-rss
./gator browse 5 Gator - rss feed aggregator
