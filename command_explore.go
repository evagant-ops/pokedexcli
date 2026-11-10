package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	println("")

	resp, err := cfg.client.GetLocationArea(cfg.nextUrl, args[0])
	if err != nil {
		return err
	}

	for _, l := range resp.PokemonEncounters {
		fmt.Println(l.Pokemon.Name)
	}
	return nil
}
