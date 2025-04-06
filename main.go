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
	Next     string
	Previous string
	Cache    *pokecache.Cache
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
		Cache: cache,
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
	return nil
}
