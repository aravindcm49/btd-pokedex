package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/pokeapi"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Pokedex > ") // Println will add back the final '\n'
	var config Config
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

	return nil
}

func commandMapb(config *Config) error {

	var prevURL string
	if config.Previous != "" {
		prevURL = config.Previous
	} else {
		fmt.Printf("you're on the first page\n")
	}

	res, err := pokeapi.GetLocationAreasList(prevURL)
	if err != nil {
		return err
	}

	config.Next = res.Next

	if res.Previous != nil {
		prevURL, ok := res.Previous.(string)
		if ok {
			config.Previous = prevURL
		}
	} else {
		config.Previous = ""
	}

	for _, la := range res.Results {
		fmt.Println(la.Name)
	}

	return nil
}

func commandConfig(config *Config) error {
	fmt.Printf("\nCONFIG\n------\nNext - %v\nPrevious - %v\n\n", config.Next, config.Previous)
	return nil
}
