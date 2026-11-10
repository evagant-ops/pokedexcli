package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	name := args[0]
	p, ok := cfg.pokedex[name]
	if !ok {
		fmt.Printf("%s is not in your Pokedex\n", name)
		return nil
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)

	fmt.Println("Stats: ")
	for _, s := range p.Stats {
		fmt.Printf("\t%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types: ")
	for _, s := range p.Types {
		fmt.Printf("\t%s\n", s.Type.Name)
	}

	return nil
}
