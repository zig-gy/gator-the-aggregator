package main

import (
	"fmt"
	"os"

	"github.com/zig-gy/gator-the-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	_ = state{&cfg}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)

	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("error: gator needs a command to run")
	}
	
}	