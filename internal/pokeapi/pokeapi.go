package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ---------- Data Structures ----------

// LocationArea represents detailed information about a location area.
type LocationArea struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// PokemonEncounter represents a Pokémon encounter in a location area.
type PokemonEncounter struct {
	Pokemon NamedAPIResource `json:"pokemon"`
}

// NamedAPIResource represents a resource with a name and URL.
type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// LocationAreaResponse represents the response structure from the PokeAPI for location areas list.
type LocationAreaResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []LocationAreaRef `json:"results"`
}

// LocationAreaRef represents a location area in the list response.
type LocationAreaRef struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Pokemon represents the data structure for a Pokémon.
type Pokemon struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	BaseXP int    `json:"base_experience"`
	// Add other relevant fields as needed.
}

// ---------- Functions ----------

// FetchPokemon fetches data for a given Pokémon name.
func FetchPokemon(name string) (*Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", strings.ToLower(name))
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching Pokémon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return &pokemon, nil
}

// GetPokeLocations fetches location areas from the given URL.
func GetPokeLocations(url string) (areas []LocationAreaRef, nextUrl string, previousUrl string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", "", fmt.Errorf("error fetching location areas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", "", fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var data LocationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, "", "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	return data.Results, data.Next, data.Previous, nil
}

// FetchLocationArea fetches details of a location area by name.
func FetchLocationArea(areaName string) (*LocationArea, error) {
	// Replace spaces with hyphens and make it lowercase for the URL
	sanitizedAreaName := strings.ToLower(strings.ReplaceAll(areaName, " ", "-"))
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", sanitizedAreaName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching location area: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var locationArea LocationArea
	if err := json.NewDecoder(resp.Body).Decode(&locationArea); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return &locationArea, nil
}
