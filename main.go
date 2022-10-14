package main

import (
	"fmt"
	"log"
	"pokemon_ascii/clear_screen"
	"pokemon_ascii/get_data"
	"pokemon_ascii/print_color"
	"pokemon_ascii/read_files"
)

// Filenames containing important data such as configuration and info for every Pokemon
// Info obtained from PokeAPI (https://pokeapi.co/docs/v2 && https://github.com/mtslzr/pokeapi-go)
var (
	CONFIG_FILE  string = "/data/config.json"
	DATA_POKEMON string = "/data/pokemon.json"
)

// Simple function that checks if there is an error so the program terminates, printing the error
func check_error(err error) {
	if err != nil {
		log.Fatalf("ups! something happened: %v", err)
	}
}

// MAIN
func main() {

	// Clear the screen before printing
	err_clear := clear_screen.CallClear()
	check_error(err_clear)

	// Get path to executable so we can easily locate some useful files for the program
	path, err_path := get_data.Get_exec_path()
	check_error(err_path)

	// Read config JSON file that contains selection criteria
	config, err_config := get_data.Get_config(path + CONFIG_FILE)
	check_error(err_config)

	// Get data from JSON file containing Pokemon data such as names, IDs, types, among others using criteria given in config file
	pokemons, pokemon_styles, err_getdata := get_data.Choose_pokemon(config, path, path+DATA_POKEMON)
	check_error(err_getdata)

	// Get colors based on Pokemon types (default); or get Pokemon colors given specific ones (custom)
	colors, err_colors := get_data.Get_Pokemon_Colors(config, pokemons)
	check_error(err_colors)

	// Read files with ASCII arts and returns a string where each ASCII art is at the side of each other
	// so it can be printed horizontally in screen. Also, returns the max width of each ASCII art, so the
	// program knows when to print in a certain color in the next function
	ascii_art, width_ascii, err_ascii := read_files.Unify_files(path, pokemons)
	check_error(err_ascii)

	// Change color from ascii art just obtained and prints each Pokemon with the color chosen by config file
	ascii_in_color := print_color.Print_in_colors(width_ascii, pokemon_styles, ascii_art, colors)

	// Finally, print the result
	fmt.Println(ascii_in_color)
}
