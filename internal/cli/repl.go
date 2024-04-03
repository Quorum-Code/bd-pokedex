package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Quorum-Code/bd-pokedex/internal/cli/config"
)

type responseData struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func Run() {
	cfg := newCfg()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		args := strings.Split(scanner.Text(), " ")
		cmd := strings.ToLower(args[0])

		cmdCallback, ok := cfg.Commands[cmd]
		if !ok {
			fmt.Println("invalid command")
		} else {
			err := cmdCallback.Callback(cfg, args[1:])
			if err == nil {
				continue
			}
			if err.Error() == "exit command" {
				fmt.Println("exiting program...")
				return
			} else {
				fmt.Println(err)
			}
		}
	}
}

func newCfg() *config.Clicfg {
	url := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	return &config.Clicfg{
		Cache:         config.NewCache(time.Second * 3),
		Commands:      buildCommands(),
		CaughtPokemon: []string{},
		MapLast:       &url,
		MapNext:       nil,
		MapPrev:       nil,
	}
}

func buildCommands() map[string]config.CliCommand {
	return map[string]config.CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays 20 map locations.",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays 20 map locations.",
			Callback:    commandMapB,
		},
		"entries": {
			Name:        "entries",
			Description: "Displays urls of all cached entries",
			Callback:    commandEntries,
		},
		"explore": {
			Name:        "explore",
			Description: "Displays Pokemon available in the area",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempts to catch a pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Reveals information about a pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Lists the pokemon the user has caught",
			Callback:    commandPokedex,
		},
	}
}

func commandHelp(cfg *config.Clicfg, args []string) error {
	fmt.Print("Usage:\n\n")

	for key := range cfg.Commands {
		fmt.Printf("%s: %s\n", cfg.Commands[key].Name, cfg.Commands[key].Description)
	}
	fmt.Println()

	return nil
}

func commandExit(cfg *config.Clicfg, args []string) error {
	return errors.New("exit command")
}

func commandEntries(cfg *config.Clicfg, args []string) error {
	for e := range cfg.Cache.Entries {
		fmt.Println(e)
	}

	return nil
}

func commandMap(cfg *config.Clicfg, args []string) error {
	return subCommandMap(cfg, cfg.MapNext)
}

func commandMapB(cfg *config.Clicfg, args []string) error {
	return subCommandMap(cfg, cfg.MapPrev)
}

func subCommandMap(cfg *config.Clicfg, url *string) error {
	if url == nil {
		url = cfg.MapLast
	}

	body, err := cfg.Cache.Get(*url)
	if err != nil {
		resp, err := http.Get(*url)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode > 299 {
			return errors.New("response failed")
		}
		if err != nil {
			return err
		}
		defer cfg.Cache.Add(*url, body)
	}

	respData := responseData{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	cfg.MapLast = url
	cfg.MapNext = respData.Next
	cfg.MapPrev = respData.Previous

	for i := range respData.Results {
		fmt.Printf("%s\n", respData.Results[i].Name)
	}

	return nil
}
