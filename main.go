package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	pokeapi "github.com/adamararcane/CLIPokedex/internal/pokeapi"
	pokecache "github.com/adamararcane/CLIPokedex/internal/pokecache"
)

// ---------- Data Structures ----------

// Struct to navigate world
type Config struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// Structures CLI commands
type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

// Init
var config Config
var commands map[string]cliCommand
var cache *pokecache.Cache
var Pokedex = make(map[string]pokeapi.Pokemon)

func init() {
	// Initialize the cache with a cleanup interval of 60 seconds
	cache = pokecache.NewCache(60 * time.Second)
}

func main() {
	if err := loadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Init commands
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokédex",
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
		"catch": {
			name:        "catch",
			description: "Attempt to capture a Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get detailed stats about a Pokémon",
			callback:    commandInspect,
		},
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the CLI Pokédex! Type 'help' for available commands.")

	for {
		fmt.Print("Pokédex > ")

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
				fmt.Println("Thanks for using the CLI Pokédex!")
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

// ---------- Functions ----------

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

func commandCatch(args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: catch <pokemon_name>")
		return nil
	}

	pokemonName := args[0]
	pokeball := 100

	// Cache pokemon for failed capture
	cacheKey := "pokemon:" + pokemonName

	var targetPokemon pokeapi.Pokemon

	if cachedData, found := cache.Get(cacheKey); found {
		err := json.Unmarshal(cachedData, &targetPokemon)
		if err != nil {
			fmt.Println("Error unmarshaling cached data:", err)
		} else {
			fmt.Printf("(from cache) Throwing a pokeball at %s.", targetPokemon.Name)
			catch := rand.Intn(targetPokemon.BaseXP)
			time.Sleep(time.Second)
			fmt.Print(".")
			time.Sleep(time.Second)
			fmt.Print(".\n")
			if catch < pokeball {
				fmt.Printf("%s was caught!\n", targetPokemon.Name)
				fmt.Printf("Adding %s to the Pokédex!\n", targetPokemon.Name)
				Pokedex[targetPokemon.Name] = targetPokemon

			} else {
				fmt.Printf("%s escaped!\n", targetPokemon.Name)
			}

			return nil
		}
	}

	// Pokemon not in cache; fetching from API
	pokemonData, err := pokeapi.FetchPokemon(pokemonName)
	if err != nil {
		fmt.Println("Error fetching pokemon:", err)
	}

	// Caching data for future catch attempts
	dataToCache, err := json.Marshal(pokemonData)
	if err != nil {
		fmt.Println("Error marshaling data for cache:", err)
	} else {
		cache.Add(cacheKey, dataToCache)
	}

	// Display Pokemon name and baseXP
	fmt.Printf("Throwing a pokeball at %s.", pokemonData.Name)
	catch := rand.Intn(pokemonData.BaseXP)
	time.Sleep(time.Second)
	fmt.Print(".")
	time.Sleep(time.Second)
	fmt.Print(".\n")
	if catch < pokeball {
		fmt.Printf("%s was caught!\n", pokemonData.Name)
		fmt.Printf("Adding %s to the Pokédex!\n", pokemonData.Name)
		Pokedex[pokemonData.Name] = *pokemonData
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
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

func commandInspect(args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: inspect <pokemon_name>")
		return nil
	}

	pokemonName := args[0]

	if _, found := Pokedex[pokemonName]; !found {
		fmt.Printf("You have not caught %s!\n", pokemonName)
		return nil
	}

	fmt.Printf("Name: %s (ID: %d)\n", pokemonName, Pokedex[pokemonName].ID)
	fmt.Printf("Height: %d\n", Pokedex[pokemonName].Height)
	fmt.Printf("Weight: %d\n", Pokedex[pokemonName].Weight)
	fmt.Println("Types:")
	for _, TypeEntry := range Pokedex[pokemonName].Type {
		fmt.Printf("  - %s\n", TypeEntry.Type.Name)
	}
	fmt.Println("Stats:")
	for _, StatEntry := range Pokedex[pokemonName].Stats {
		fmt.Printf("  - %s: %d\n", StatEntry.Stat.Name, StatEntry.BaseStat)
	}

	return nil
}

// ---------- CONFIG FUNCTIONS ----------

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
