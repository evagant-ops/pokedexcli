package main

import (
	"errors"
	"fmt"
	"math/rand"
)


const catchDifficulty = 60

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}
	name := args[0]

	if _, ok := cfg.pokedex[name]; ok {
		fmt.Printf("%s is already in your Pokedex\n", name)
		return nil
	}

	pokemon, err := cfg.client.GetPokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if pokemon.BaseExperience > 0 && rand.Intn(pokemon.BaseExperience) >= catchDifficulty {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)
	cfg.pokedex[pokemon.Name] = pokemon
	return nil
}
