package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	pokeapi "github.com/adamararcane/CLIPokedex/internal/pokeapi"
	pokecache "github.com/adamararcane/CLIPokedex/internal/pokecache"
)

type Config struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

var config Config
var commands map[string]cliCommand
var cache *pokecache.Cache

func init() {
	// Initialize the cache with a cleanup interval of 60 seconds
	cache = pokecache.NewCache(60 * time.Second)
}

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
		"explore": {
			name:        "explore",
			description: "Explore a location area to see Pokémon",
			callback:    commandExplore,
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
		tokens := strings.Fields(input)
		if len(tokens) == 0 {
			continue // No input entered
		}

		commandName := tokens[0]
		args := tokens[1:]

		if cmd, exists := commands[commandName]; exists {
			if err := cmd.callback(args); err != nil {
				fmt.Println("Error:", err)
			}
			if commandName == "exit" {
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
	callback    func(args []string) error
}

func commandHelp(args []string) error {
	fmt.Println("===== Commands =====")
	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Println("\nUsage:")
	fmt.Println("  explore <area_name> - Explore a location area to see Pokémon")
	// Add usage examples if necessary
	return nil
}

func commandExit(args []string) error {
	fmt.Println("Exiting the Pokedex")
	return nil
}

func commandMap(args []string) error {
	url := config.Next
	if url == "" {
		fmt.Println("No next page available.")
		return nil
	}

	var areas []pokeapi.LocationAreaRef
	var nextUrl, previousUrl string

	cacheKey := url

	if cachedData, found := cache.Get(cacheKey); found {
		err := json.Unmarshal(cachedData, &areas)
		if err != nil {
			fmt.Println("Error unmarshaling cached data:", err)
		} else {
			nextUrl = config.Next
			previousUrl = config.Previous
		}
	}

	if areas == nil || len(areas) == 0 {
		var err error
		areas, nextUrl, previousUrl, err = pokeapi.GetPokeLocations(url)
		if err != nil {
			return err
		}

		dataToCache, err := json.Marshal(areas)
		if err != nil {
			fmt.Println("Error marshaling data for cache:", err)
		} else {
			cache.Add(cacheKey, dataToCache)
		}
	}

	config.Next = nextUrl
	config.Previous = previousUrl
	if err := saveConfig(); err != nil {
		fmt.Println("Warning: Could not save config.", err)
	}

	fmt.Println("Location Areas:")
	for i, area := range areas {
		fmt.Printf("%d. %s\n", i+1, area.Name)
	}

	return nil
}

func commandMapb(args []string) error {
	url := config.Previous
	if url == "" {
		fmt.Println("No previous page available.")
		return nil
	}

	var areas []pokeapi.LocationAreaRef
	var nextUrl, previousUrl string

	cacheKey := url

	if cachedData, found := cache.Get(cacheKey); found {
		err := json.Unmarshal(cachedData, &areas)
		if err != nil {
			fmt.Println("Error unmarshaling cached data:", err)
		} else {
			nextUrl = config.Next
			previousUrl = config.Previous
		}
	}

	if areas == nil || len(areas) == 0 {
		var err error
		areas, nextUrl, previousUrl, err = pokeapi.GetPokeLocations(url)
		if err != nil {
			return err
		}

		dataToCache, err := json.Marshal(areas)
		if err != nil {
			fmt.Println("Error marshaling data for cache:", err)
		} else {
			cache.Add(cacheKey, dataToCache)
		}
	}

	config.Next = nextUrl
	config.Previous = previousUrl
	if err := saveConfig(); err != nil {
		fmt.Println("Warning: Could not save config.", err)
	}

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

func commandExplore(args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: explore <area_name>")
		return nil
	}

	areaName := args[0]

	// Use cache key based on area name
	cacheKey := "explore:" + areaName

	var pokemonNames []string

	// Check if data is in cache
	if cachedData, found := cache.Get(cacheKey); found {
		err := json.Unmarshal(cachedData, &pokemonNames)
		if err != nil {
			fmt.Println("Error unmarshaling cached data:", err)
			// Proceed to fetch from API
		} else {
			fmt.Printf("Pokémon in %s (from cache):\n", areaName)
			for _, name := range pokemonNames {
				fmt.Println("- " + name)
			}
			return nil
		}
	}

	// Data not in cache; fetch from API
	locationArea, err := pokeapi.FetchLocationArea(areaName)
	if err != nil {
		fmt.Println("Error fetching location area:", err)
		return nil
	}

	// Extract Pokémon names
	for _, encounter := range locationArea.PokemonEncounters {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}

	// Cache the Pokémon names
	dataToCache, err := json.Marshal(pokemonNames)
	if err != nil {
		fmt.Println("Error marshaling data for cache:", err)
	} else {
		cache.Add(cacheKey, dataToCache)
	}

	// Display Pokémon names
	fmt.Printf("Pokémon in %s:\n", areaName)
	for _, name := range pokemonNames {
		fmt.Println("- " + name)
	}

	return nil
}
