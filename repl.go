package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/evagant-ops/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}
type config struct {
	pokedex  map[string]pokeapi.Pokemon
	client   *pokeapi.Client
	commands map[string]cliCommand
	nextUrl  string
	prevUrl  string
}

func cleanInput(input string) []string {
	words := strings.Fields(input)

	for i, word := range words {
		if word == strings.ToUpper(word) {
			words[i] = strings.ToLower(word)
		}
	}

	return words
}

func startRepl(cfg *config) {
	const pokedex = "Pokedex >"
	fmt.Print(pokedex)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		words := cleanInput(text)
		if len(words) > 0 {
			commands := getCommands(cfg)
			if cmd, ok := commands[words[0]]; ok {
				err := cmd.callback(cfg, words[1:]...)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Printf("Your command was: %s\n", strings.ToLower(words[0]))
			}
		}
		fmt.Print(pokedex)
	}
}

func getCommands(cfg *config) map[string]cliCommand {
	return cfg.commands
}
