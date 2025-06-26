package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/zig-gy/gator-the-aggregator/internal/config"
	"github.com/zig-gy/gator-the-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbURL := cfg.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	s := state{
		cfg: &cfg,
		db: dbQueries,
	}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("error: gator needs a command to run")
		os.Exit(1)
	}

	cmd := command{
		name: arguments[1],
		arguments: arguments[2:],
	}

	for name := range cmds.cmds {
		if cmd.name == name {
			if err := cmds.run(&s, cmd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
	
	fmt.Printf("error: command %s not found\n", cmd.name)
	os.Exit(1)
}	
