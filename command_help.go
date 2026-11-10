package main

import "fmt"

func commandHelp(cfg *config, args ...string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")
	for _, cmd := range getCommands(cfg) {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
