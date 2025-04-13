package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next               string
	Previous           string
	AddtionalArguments string
	Pokedex            map[string]pokeapi.Pokemon
	Cache              *pokecache.Cache
}

var helpmap map[string]cliCommand

func init() {
	helpmap = map[string]cliCommand{
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
			description: "List of pokemon locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back in the pokemon locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore and find pokemon in a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "try to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "display your pokedex",
			callback:    commandPokedex,
		},
		"config": {
			name:        "config",
			description: "Print out the config",
			callback:    commandConfig,
		},
	}
}

func main() {
	cache := pokecache.NewCache(5 * time.Minute)
	config := Config{
		Cache:   cache,
		Pokedex: make(map[string]pokeapi.Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Pokedex > ") // Println will add back the final '\n'
	for scanner.Scan() {
		textInput := scanner.Text()
		currentCommand := cleanInput(textInput)[0]
		res, ok := helpmap[currentCommand]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			if len(cleanInput(textInput)) > 1 {
				config.AddtionalArguments = strings.Join(cleanInput(textInput)[1:], " ")
			} else {
				config.AddtionalArguments = ""
			}
			res.callback(&config)
		}
		fmt.Printf("Pokedex > ") // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func cleanInput(text string) []string {
	res := strings.ToLower(text)
	return strings.Fields(res)
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, command := range helpmap {
		fmt.Printf("%v : %v\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *Config) error {
	var nextURL string
	if config.Next != "" {
		nextURL = config.Next
	} else {
		nextURL = "https://pokeapi.co/api/v2/location-area/"
	}

	cachedData, ok := config.Cache.Get(nextURL)
	if ok {
		// fmt.Println("using cached data")
		var res pokeapi.LocationAreaList
		err := json.Unmarshal(cachedData, &res)
		if err != nil {
			return err
		}
		config.Next = res.Next
		config.Previous = res.Previous.(string)
		for _, la := range res.Results {
			fmt.Println(la.Name)
		}
		return nil
	}

	res, err := pokeapi.GetLocationAreasList(nextURL)
	if err != nil {
		return err
	}

	config.Next = res.Next

	if res.Previous != nil {
		config.Previous = res.Previous.(string)
	} else {
		config.Previous = ""
	}

	for _, la := range res.Results {
		fmt.Println(la.Name)
	}

	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	config.Cache.Add(nextURL, data)

	return nil
}

func commandMapb(config *Config) error {

	var prevURL string
	if config.Previous != "" {
		prevURL = config.Previous
	} else {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	cachedData, ok := config.Cache.Get(prevURL)
	if ok {
		// fmt.Println("using cached data")
		var res pokeapi.LocationAreaList
		err := json.Unmarshal(cachedData, &res)
		if err != nil {
			return err
		}
		config.Next = res.Next
		if res.Previous != nil {
			config.Previous = res.Previous.(string)
		} else {
			config.Previous = ""
		}
		for _, la := range res.Results {
			fmt.Println(la.Name)
		}
		return nil
	}

	res, err := pokeapi.GetLocationAreasList(prevURL)
	if err != nil {
		return err
	}

	config.Next = res.Next

	if res.Previous != nil {
		config.Previous = res.Previous.(string)
	} else {
		config.Previous = ""
	}

	for _, la := range res.Results {
		fmt.Println(la.Name)
	}

	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	config.Cache.Add(prevURL, data)

	return nil
}

func commandConfig(config *Config) error {
	fmt.Printf("\nCONFIG\n------\nNext - %v\nPrevious - %v\n\n", config.Next, config.Previous)
	fmt.Printf("Pokedex:\n")
	for _, v := range config.Pokedex {
		fmt.Printf("- %v\n", v.Name)
	}
	return nil
}

func commandExplore(config *Config) error {
	secondInput := config.AddtionalArguments
	if secondInput == "" {
		fmt.Printf("Please provide a location name\n")
		return nil
	}
	res, err := pokeapi.GetPokemonsForLA(secondInput)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %v...\n", secondInput)
	fmt.Println("Found Pokemon:")
	for _, pe := range res {
		fmt.Printf("- %v\n", pe.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config) error {
	secondInput := config.AddtionalArguments
	if secondInput == "" {
		fmt.Printf("Please provide a pokemon name\n")
		return nil
	}

	isPokeCaught, res, err := pokeapi.CatchPokemon(secondInput)
	if err != nil {
		fmt.Printf("Error with api pokemon: %v\n", err)
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", secondInput)
	if isPokeCaught {
		fmt.Printf("%v was caught!\n", res.Name)
		config.Pokedex[res.Name] = res
	} else {
		fmt.Printf("%v escaped!\n", res.Name)
	}

	return nil
}

func commandInspect(config *Config) error {
	secondInput := config.AddtionalArguments
	if secondInput == "" {
		fmt.Printf("Please provide a pokemon name\n")
		return nil
	}
	var res pokeapi.Pokemon
	pokeInMap, ok := config.Pokedex[secondInput]
	if ok {
		res = pokeInMap
	} else {
		poke, err := pokeapi.GetPokemon(secondInput)
		if err != nil {
			fmt.Printf("Error with api pokemon: %v\n", err)
		}
		res = poke
	}

	fmt.Printf("Name: %s\nHeight: %v\nWeight: %v\n", res.Name, res.Height, res.Weight)
	fmt.Println("Stats:")
	for _, s := range res.Stats {
		fmt.Printf("  - %v: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range res.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(config *Config) error {
	fmt.Printf("Your Pokedex:\n")
	for _, v := range config.Pokedex {
		fmt.Printf("  - %v\n", v.Name)
	}
	return nil
}
