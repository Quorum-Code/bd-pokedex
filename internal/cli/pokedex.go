package cli

import (
	"fmt"

	"github.com/Quorum-Code/bd-pokedex/internal/cli/config"
)

func commandPokedex(cfg *config.Clicfg, args []string) error {
	if len(cfg.CaughtPokemon) >= 0 {
		fmt.Print("You have no pokemon...\n")
		return nil
	}

	fmt.Print("Your pokemon\n")
	for i := range cfg.CaughtPokemon {
		fmt.Printf("  - %s\n", cfg.CaughtPokemon[i])
	}

	return nil
}
