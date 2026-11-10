package main

import (
	"fmt"
)

func commandMap(cfg *config, args ...string) error {
	resp, err := cfg.client.GetLocationAreas(cfg.nextUrl)
	if err != nil {
		return err
	}

	for _, a := range resp.Results {
		fmt.Println(a.Name)
	}
	cfg.prevUrl = resp.Previous
	cfg.nextUrl = resp.Next

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevUrl == "" {
		fmt.Println("You are on the first page")
		return nil
	}

	resp, err := cfg.client.GetLocationAreas(cfg.prevUrl)
	if err != nil {
		return err
	}

	for _, a := range resp.Results {
		fmt.Println(a.Name)
	}

	cfg.prevUrl = resp.Previous
	cfg.nextUrl = resp.Next
	return nil
}
