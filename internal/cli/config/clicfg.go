package config

import "slices"

type Clicfg struct {
	Cache         Cache
	Commands      map[string]CliCommand
	CaughtPokemon []string
	MapLast       *string
	MapNext       *string
	MapPrev       *string
}

func NewClicfg() *Clicfg {
	return &Clicfg{}
}

func (c *Clicfg) AddPokemon(pokemon string) {
	if !slices.Contains(c.CaughtPokemon, pokemon) {
		c.CaughtPokemon = append(c.CaughtPokemon, pokemon)
	}
}
