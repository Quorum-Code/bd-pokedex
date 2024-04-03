package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Quorum-Code/bd-pokedex/internal/cli/config"
)

func commandInspect(cfg *config.Clicfg, args []string) error {
	if len(args) <= 0 {
		return errors.New("no pokemon given as argument")
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	body, err := cfg.Cache.Get(url)
	if err != nil {
		return err
	}

	respData := pokemonData{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", respData.Name)
	fmt.Printf("Height: %d\n", respData.Height)
	fmt.Printf("Weight: %d\n", respData.Weight)

	fmt.Printf("Stats:\n")
	for i := range respData.Stats {
		fmt.Printf("  -%s: %d\n", respData.Stats[i].Stat.Name, respData.Stats[i].BaseStat)
	}

	fmt.Printf("Types:\n")
	for i := range respData.Types {
		fmt.Printf("  - %s\n", respData.Types[i].Type.Name)
	}

	return nil
}
