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
)

type cliCommand struct {
	name        string
	description string
	callback    func(*clicfg) error
}

type clicfg struct {
	cache    Cache
	commands map[string]cliCommand
	mapLast  *string
	mapNext  *string
	mapPrev  *string
}

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
		text := scanner.Text()
		text = strings.ToLower(text)

		cmd, ok := cfg.commands[text]
		if !ok {
			fmt.Println("invalid command")
		} else {
			err := cmd.callback(cfg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func newCfg() *clicfg {
	url := "https://pokeapi.co/api/v2/location/?offset=0&limit=20"

	return &clicfg{
		NewCache(time.Second * 3),
		buildCommands(),
		&url,
		nil,
		nil,
	}
}

func buildCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays 20 map locations.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 map locations.",
			callback:    commandMapB,
		},
		"entries": {
			name:        "entries",
			description: "Displays urls of all cached entries",
			callback:    commandEntries,
		},
	}
}

func commandHelp(cfg *clicfg) error {
	fmt.Print("Usage:\n\n")

	for key := range cfg.commands {
		fmt.Printf("%s: %s\n", cfg.commands[key].name, cfg.commands[key].description)
	}
	fmt.Println()

	return nil
}

func commandExit(cfg *clicfg) error {
	return errors.New("exit command")
}

func commandEntries(cfg *clicfg) error {
	for e := range cfg.cache.entries {
		fmt.Println(e)
	}

	return nil
}

func commandMap(cfg *clicfg) error {
	return subCommandMap(cfg, cfg.mapNext)
}

func commandMapB(cfg *clicfg) error {
	return subCommandMap(cfg, cfg.mapPrev)
}

func subCommandMap(cfg *clicfg, url *string) error {
	if url == nil {
		url = cfg.mapLast
	}

	body, err := cfg.cache.Get(*url)
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
		defer cfg.cache.Add(*url, body)
	}

	respData := responseData{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	cfg.mapLast = url
	cfg.mapNext = respData.Next
	cfg.mapPrev = respData.Previous

	for i := range respData.Results {
		fmt.Printf("%s\n", respData.Results[i].Name)
	}

	return nil
}
