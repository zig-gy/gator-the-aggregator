# gator the aggreGATOR

boot.dev course on SQL and Go programming.

gator is a way to aggregate your favourite RSS feeds in your terminal.

## Requirements

- PostgreSQL
- Go

## How to build

1. Clone the repo
2. Create a database for the program to use
3. Create a file called `.gatorconfig.json` in your home directory with this format
```
{
    "db_url":               "connection string here",
    "current_user_name":    "leave blank"
}
```
4. Using goose (`go install github.com/pressly/goose/v3/cmd/goose@latest` if you don't have it) execute an up migration from the `sql/schema` directory
5. From the project directory execute `sqlc generate` (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` if you don't have it)
6. Execute `go install`!

## How to use

In the terminal execute `gator-the-aggregator` with a command to interact with RSS feeds.

This is the list of supported commands:

- **login**: Change the current user if it exists (needs 1 argument to run, the user's name)
- **register**: Create a new user and login as that user (needs 1 argument to run, the user's name)
- **reset**: Delete every user and the data associated with it
- **users**: List every user
- **agg**: Start the process of fetching the feeds waiting for the specified time (needs 1 argument to , the time to wait between fetches, eg: 1m = 1 minute)
- **addfeed**: Add a feed to the database and follow it (needs 2 arguments to run, name of the feed and the feed's url)
- **feeds**: Show all the feeds
- **follow**: Follow an already added feed (needs 1 argument, the feed's url)
- **following**: Show all the feed the current user follows
- **unfollow**: Unfollow a feed (needs 1 argument, the feed's url)
- **browse**: Show your latest posts (accepts 1 argument, the number of posts to show, defaults to 2)