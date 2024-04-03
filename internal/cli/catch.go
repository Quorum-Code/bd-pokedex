package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/Quorum-Code/bd-pokedex/internal/cli/config"
)

type pokemonData struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	BaseExperience int    `json:"base_experience,omitempty"`
	Height         int    `json:"height,omitempty"`
	IsDefault      bool   `json:"is_default,omitempty"`
	Order          int    `json:"order,omitempty"`
	Weight         int    `json:"weight,omitempty"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden,omitempty"`
		Slot     int  `json:"slot,omitempty"`
		Ability  struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"ability,omitempty"`
	} `json:"abilities,omitempty"`
	Forms []struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"forms,omitempty"`
	GameIndices []struct {
		GameIndex int `json:"game_index,omitempty"`
		Version   struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"version,omitempty"`
	} `json:"game_indices,omitempty"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"item,omitempty"`
		VersionDetails []struct {
			Rarity  int `json:"rarity,omitempty"`
			Version struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"held_items,omitempty"`
	LocationAreaEncounters string `json:"location_area_encounters,omitempty"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"move,omitempty"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at,omitempty"`
			VersionGroup   struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version_group,omitempty"`
			MoveLearnMethod struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"move_learn_method,omitempty"`
		} `json:"version_group_details,omitempty"`
	} `json:"moves,omitempty"`
	Species struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"species,omitempty"`
	Sprites struct {
		BackDefault      string `json:"back_default,omitempty"`
		BackFemale       any    `json:"back_female,omitempty"`
		BackShiny        string `json:"back_shiny,omitempty"`
		BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
		FrontDefault     string `json:"front_default,omitempty"`
		FrontFemale      any    `json:"front_female,omitempty"`
		FrontShiny       string `json:"front_shiny,omitempty"`
		FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default,omitempty"`
				FrontFemale  any    `json:"front_female,omitempty"`
			} `json:"dream_world,omitempty"`
			Home struct {
				FrontDefault     string `json:"front_default,omitempty"`
				FrontFemale      any    `json:"front_female,omitempty"`
				FrontShiny       string `json:"front_shiny,omitempty"`
				FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
			} `json:"home,omitempty"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default,omitempty"`
				FrontShiny   string `json:"front_shiny,omitempty"`
			} `json:"official-artwork,omitempty"`
			Showdown struct {
				BackDefault      string `json:"back_default,omitempty"`
				BackFemale       any    `json:"back_female,omitempty"`
				BackShiny        string `json:"back_shiny,omitempty"`
				BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
				FrontDefault     string `json:"front_default,omitempty"`
				FrontFemale      any    `json:"front_female,omitempty"`
				FrontShiny       string `json:"front_shiny,omitempty"`
				FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
			} `json:"showdown,omitempty"`
		} `json:"other,omitempty"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackGray     string `json:"back_gray,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontGray    string `json:"front_gray,omitempty"`
				} `json:"red-blue,omitempty"`
				Yellow struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackGray     string `json:"back_gray,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontGray    string `json:"front_gray,omitempty"`
				} `json:"yellow,omitempty"`
			} `json:"generation-i,omitempty"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"crystal,omitempty"`
				Gold struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"gold,omitempty"`
				Silver struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"silver,omitempty"`
			} `json:"generation-ii,omitempty"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"emerald,omitempty"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"firered-leafgreen,omitempty"`
				RubySapphire struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"ruby-sapphire,omitempty"`
			} `json:"generation-iii,omitempty"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"diamond-pearl,omitempty"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"heartgold-soulsilver,omitempty"`
				Platinum struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"platinum,omitempty"`
			} `json:"generation-iv,omitempty"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default,omitempty"`
						BackFemale       any    `json:"back_female,omitempty"`
						BackShiny        string `json:"back_shiny,omitempty"`
						BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
						FrontDefault     string `json:"front_default,omitempty"`
						FrontFemale      any    `json:"front_female,omitempty"`
						FrontShiny       string `json:"front_shiny,omitempty"`
						FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
					} `json:"animated,omitempty"`
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"black-white,omitempty"`
			} `json:"generation-v,omitempty"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"omegaruby-alphasapphire,omitempty"`
				XY struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"x-y,omitempty"`
			} `json:"generation-vi,omitempty"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontFemale  any    `json:"front_female,omitempty"`
				} `json:"icons,omitempty"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"ultra-sun-ultra-moon,omitempty"`
			} `json:"generation-vii,omitempty"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontFemale  any    `json:"front_female,omitempty"`
				} `json:"icons,omitempty"`
			} `json:"generation-viii,omitempty"`
		} `json:"versions,omitempty"`
	} `json:"sprites,omitempty"`
	Cries struct {
		Latest string `json:"latest,omitempty"`
		Legacy string `json:"legacy,omitempty"`
	} `json:"cries,omitempty"`
	Stats []struct {
		BaseStat int `json:"base_stat,omitempty"`
		Effort   int `json:"effort,omitempty"`
		Stat     struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"stat,omitempty"`
	} `json:"stats,omitempty"`
	Types []struct {
		Slot int `json:"slot,omitempty"`
		Type struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"type,omitempty"`
	} `json:"types,omitempty"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"generation,omitempty"`
		Types []struct {
			Slot int `json:"slot,omitempty"`
			Type struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"type,omitempty"`
		} `json:"types,omitempty"`
	} `json:"past_types,omitempty"`
}

func commandCatch(cfg *config.Clicfg, args []string) error {
	if len(args) <= 0 {
		return errors.New("no pokemon argument given")
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	body, err := cfg.Cache.Get(url)
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
		defer cfg.Cache.Add(url, body)
	}

	respData := pokemonData{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	val := rand.Float64()

	if val >= 25.0/float64(respData.BaseExperience) {
		fmt.Printf("Caught %s!\n", args[0])
	} else {
		fmt.Printf("%s escaped!\n", args[0])
	}

	return nil
}
