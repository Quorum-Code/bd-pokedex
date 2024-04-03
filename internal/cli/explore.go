package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type exploreData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name *string `json:"name"`
			URL  *string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name *string `json:"name"`
				URL  *string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name *string `json:"name"`
		URL  *string `json:"url"`
	} `json:"location"`
	Name  *string `json:"name"`
	Names []struct {
		Language struct {
			Name *string `json:"name"`
			URL  *string `json:"url"`
		} `json:"language"`
		Name *string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name *string `json:"name"`
			URL  *string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name *string `json:"name"`
					URL  *string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name *string `json:"name"`
				URL  *string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExplore(cfg *clicfg, args []string) error {
	if len(args) <= 0 {
		return errors.New("no location argument given")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	body, err := cfg.cache.Get(url)
	if err != nil {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode > 299 {
			return errors.New("failed response")
		}
		if err != nil {
			return nil
		}
		defer cfg.cache.Add(url, body)
	}

	respData := exploreData{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	for i := range respData.PokemonEncounters {
		fmt.Printf(" - %s\n", *respData.PokemonEncounters[i].Pokemon.Name)
	}

	return nil
}
