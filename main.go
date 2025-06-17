package main

import (
	"fmt"

	"github.com/zig-gy/gator-the-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(err)
		return
	}
	
	if err := cfg.SetUser("benjamin"); err != nil {
		fmt.Print(err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(cfg.CurrentUsername)
	fmt.Println(cfg.DbUrl)	
}