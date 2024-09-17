package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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

}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func commandHelp() error {
	fmt.Println("===== Commands =====")
	fmt.Println("Usage:")
	fmt.Printf("%s: %s\n", commands["help"].name, commands["help"].description)
	fmt.Printf("%s: %s\n", commands["exit"].name, commands["exit"].description)
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the Pokedex")
	return nil
}
