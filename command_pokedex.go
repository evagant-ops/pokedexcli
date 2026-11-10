package main

import (
	"errors"
	"fmt"
)

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.pokedex) == 0 {
		return errors.New("no pokedex found")
	}

	fmt.Println("Your Pokedex:")
	for _, p := range cfg.pokedex {
		fmt.Printf("- %s\n", p.Name)
	}
	return nil
}
