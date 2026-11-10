package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/evagant-ops/pokedexcli/internal/pokeapi"
	"github.com/evagant-ops/pokedexcli/internal/pokecache"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	defer cancel()
	cache := pokecache.NewCache(context.Background(), 5*time.Minute)
	cfg := &config{
		client:  pokeapi.NewClient(cache),
		pokedex: map[string]pokeapi.Pokemon{},
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Displays next 20 location areas",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Displays previous 20 location areas",
				callback:    commandMapb,
			},
			"explore": {
				name:        "explore",
				description: "Expores pokemons in an area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Catch pokemons",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect pokemon in pokedex",
				callback:    commandInspect,
			},
			"pokedex": {
				name:        "pokedex",
				description: "Prints your pokedex",
				callback:    commandPokedex,
			},
		}}
	go startRepl(cfg)
	<-ctx.Done()
}
