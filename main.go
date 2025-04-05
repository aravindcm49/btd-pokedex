package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Pokedex > ") // Println will add back the final '\n'
	for scanner.Scan() {
		textInput := scanner.Text()
		currentCommand := cleanInput(textInput)[0]
		res, ok := helpmap[currentCommand]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			res.callback()
		}
		// fmt.Printf("Your command was: %v\n", cleanedInput[0])

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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, command := range helpmap {
		fmt.Printf("%v : %v\n", command.name, command.description)
	}
	return nil
}
