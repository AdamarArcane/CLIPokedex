package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/adamararcane/CLIPokedex/internal"
)

type Config struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

var config Config
var commands map[string]cliCommand

func main() {
	if err := loadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	commands = map[string]cliCommand{
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
			description: "Displays next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    commandMapb,
		},
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the CLI Pokedex! Type 'help' for available commands.")

	for {
		fmt.Print("Pokedex > ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
		}

		input = strings.TrimSpace(input)

		if cmd, exists := commands[input]; exists {
			if err := cmd.callback(); err != nil {
				fmt.Println("Error: ", err)
			}
			if input == "exit" {
				fmt.Println("Thanks for using the CLI Pokedex!")
				break
			}
		} else {
			fmt.Println("Invalid command! Type 'help' for available commands.")
		}
	}

	if err := saveConfig(); err != nil {
		fmt.Println("Warning: Could not save config.", err)
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println("===== Commands =====")
	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the Pokedex")
	return nil
}

func commandMap() error {
	url := config.Next
	if url == "" {
		fmt.Println("No next page available.")
		return nil
	}

	areas, nextUrl, previousUrl, err := pokeapi.GetPokeLocations(url)
	if err != nil {
		return err
	}

	// Update config
	config.Next = nextUrl
	config.Previous = previousUrl

	if err := saveConfig(); err != nil {
		fmt.Println("Warning: Could not save config.", err)
	}

	// Display the location area names
	fmt.Println("Location Areas:")
	for i, area := range areas {
		fmt.Printf("%d. %s\n", i+1, area.Name)
	}

	return nil
}

func commandMapb() error {
	url := config.Previous
	if url == "" {
		fmt.Println("No previous page available.")
		return nil
	}

	areas, nextUrl, previousUrl, err := pokeapi.GetPokeLocations(url)
	if err != nil {
		return err
	}

	// Update config
	config.Next = nextUrl
	config.Previous = previousUrl

	if err := saveConfig(); err != nil {
		fmt.Println("Warning: Could not save config.", err)
	}

	// Display the location area names
	fmt.Println("Location Areas:")
	for i, area := range areas {
		fmt.Printf("%d. %s\n", i+1, area.Name)
	}

	return nil
}

func saveConfig() error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize default config
			config.Next = "https://pokeapi.co/api/v2/location-area/?limit=20"
			config.Previous = ""
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	return nil
}
