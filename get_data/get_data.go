package get_data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Foldername where custom pokemon choices will be alocated
const custom_foldername string = "custom"

// Define a struct for Pokemon and store their data
type Pokemons struct { //I know Pokemon plural is also Pokemon. Do not kill me
	Pokemons []Pokemon `json:"pokemon"`
}

// Pokemon attributes
type Pokemon struct {
	Name   string  `json:"name"`
	ID     int     `json:"ID"`
	Type   string  `json:"type"`
	Gen    int     `json:"gen"`
	Group  string  `json:"group"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
	Stage  string  `json:"stage"`
}

// Struct storing general config file data
type Config struct {
	ChoosePokemon string `json:"choosePokemon"`
	Group         string `json:"group"`
	Color         string `json:"color"`
	Style         string `json:"style"`
	Gen           int    `json:"gen"`
	RepeatType    bool   `json:"repeatType"`
	RepeatColor   bool   `json:"repeatColor"` //only will be used if color is set to 'random' in config file
	ShowDetails   bool   `json:"showDetails"`
	WhiteTerminal bool   `json:"whiteTerminal"`
}

// Struct for custom Pokemon files
type Custom_pokemons struct {
	Custom_pokemons []Custom_pokemon `json:"pokemon"`
}

type Custom_pokemon struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Style string `json:"style"`
}

// ANSI colors.
// This function always returns the same list with valid colors and is immutable, i.e., this is like a 'constant' list
func valid_colors() [16]string {
	return [16]string{"red", "green", "brown", "blue", "purple", "cyan", "light_gray", "dark_gray", "white",
		"light_red", "light_green", "yellow", "light_blue", "light_purple", "light_cyan", "black"}
}

// Read config file with parameters chosen by the user
func Get_config(config_file string) (Config, error) {

	var config Config

	err := read_json_default_config(&config, config_file)
	if err != nil {
		return config, err
	}
	return config, nil
}

func Choose_pokemon(config Config, path, data_pokemon string) (Pokemons, []string, error) {

	//Choose Pokemon if user uses default config file with its parameters
	if config.ChoosePokemon == "default" {
		pokemons, pokemon_styles, err := choose_pokemon_with_default_config(config, data_pokemon)
		if err != nil {
			return pokemons, pokemon_styles, err
		}
		return pokemons, pokemon_styles, nil
	}

	//If the user wants a custom config then read it from their custom config file.
	//Only valid if this custom config is located in the directory where the default config also is
	pokemons, pokemon_styles, err := choose_pokemon_with_custom_config(path, config.ChoosePokemon, data_pokemon)
	if err != nil {
		return pokemons, pokemon_styles, err
	}
	return pokemons, pokemon_styles, nil
}

// If the user wants a custom choice for Pokemons then start checking and reading from custom file
func choose_pokemon_with_custom_config(path, config_name, data_pokemon string) (Pokemons, []string, error) {
	var (
		pokemons        Pokemons
		custom_pokemons Custom_pokemons
		pokemon_styles  []string
	)

	//First check if the folder, where custom Pokemons files should be, exists
	//If it does not exist returns an error, but before quitting it creates the folder '/custom'
	//If you want to create a folder with different name change the value of global constant string 'custom_foldername'
	err := check_custom_files(path + "/" + custom_foldername)
	if err != nil {
		return pokemons, pokemon_styles, err
	}

	//Read Pokemons chosen in custom JSON file
	err_custom := read_json_custom_pokemon(&custom_pokemons, path+"/"+custom_foldername+"/"+config_name+".json")
	if err_custom != nil {
		return pokemons, pokemon_styles, err_custom
	}

	//Read file containing all Pokemon data
	all_pokemons, err := extract_Pokemon_data(data_pokemon, false)
	if err != nil {
		return pokemons, pokemon_styles, err
	}

	//Convert custom Pokemon struct data into 'Pokemon' variable struct type
	pokemons, pokemon_styles = pass_custom_to_normal_Pokemon(custom_pokemons, all_pokemons)
	if len(pokemons.Pokemons) != 3 {
		return pokemons, pokemon_styles, errors.New("number of custom Pokemon is different than 3 (unique value allowed)\n")
	}

	//If everything is ok, return Pokemon chosen by the user, their style and a 'non-error'
	return pokemons, pokemon_styles, nil
}

// If the user does not want to choose Pokemon by custom config, use 'default' config and extract its parameters
func choose_pokemon_with_default_config(config Config, data_pokemon string) (Pokemons, []string, error) {

	var (
		pokemons               Pokemons
		IDs_aready_selected    []int
		types_already_selected []string
		pokemon_styles             = []string{config.Style, config.Style, config.Style}
		counter                int = 0
	)

	//Check if the user wants a random choice for pokemons or wants a specified group of them
	if config.Group != "random" {
		if !is_in_pokegroups(config.Group) {
			fmt.Println("Warning! Pokemon group not valid. Using default value (\"random\") instead")
			*&config.Group = "random"
		}
	}

	all_pokemons, err := extract_Pokemon_data(data_pokemon, config.ShowDetails)
	if err != nil {
		return pokemons, pokemon_styles, err
	}

	//Start selecting Pokemon randomly in Pokemon JSON file
	for {
		if counter >= 3 {
			break
		}

		//First, select randomly a Pokemon and check if the Pokemon selected has not been already selected before in previous loops
		id := choose_random_pokemon_by_gen(config.Gen)
		if check_if_ID_is_already_selected(id, IDs_aready_selected) {
			continue
		} else {
			IDs_aready_selected = append(IDs_aready_selected, id)
		}

		pokemon_chosen, err := select_pokemon_by_id(id, all_pokemons)
		if err != nil {
			return pokemons, pokemon_styles, err
		}

		//Second, and if the user wants to, check that Pokemon selected belong to the same group
		if config.Group != "random" {
			if pokemon_chosen.Group != config.Group {
				continue
			}
		}

		//Finally, and if the user wants to, check there are no Pokemon types repeated
		if !config.RepeatType {
			current_type := pokemon_chosen.Type
			if is_string_already_selected(current_type, types_already_selected) {
				continue
			} else {
				types_already_selected = append(types_already_selected, strings.ToLower(current_type))
			}
		}
		pokemons.Pokemons = append(pokemons.Pokemons, pokemon_chosen)
		counter += 1
	}

	//Check if everything went fine after the loops
	if len(pokemons.Pokemons) != 3 {
		err := errors.New("pokemons obtained are more than 3; not supported")
		return pokemons, pokemon_styles, err
	}

	//If everything went fine then just return the names and filters to be printed in the next steps
	return pokemons, pokemon_styles, nil
}

// Gets data from Pokemon JSON file
func extract_Pokemon_data(filename string, show_details bool) (Pokemons, error) {

	var pokemons Pokemons

	//Read data from JSON file containing Pokemon data
	err := read_json_big_file(&pokemons, filename)

	if err != nil {
		return pokemons, err
	}

	//If the user wants to, print some details
	if show_details {
		for i := 0; i < len(pokemons.Pokemons); i++ {
			fmt.Println("Pokemon number", i+1)
			fmt.Println("Pokemon name: ", pokemons.Pokemons[i].Name)
			fmt.Println("Pokemon ID: ", strconv.Itoa(pokemons.Pokemons[i].ID))
			fmt.Println("Pokemon type: ", pokemons.Pokemons[i].Type)
			fmt.Println("Pokemon Gen: ", strconv.Itoa(pokemons.Pokemons[i].Gen))
			fmt.Println("Pokemon Group: ", pokemons.Pokemons[i].Group)
			fmt.Println("Pokemon evolution stage: ", pokemons.Pokemons[i].Stage)
			fmt.Println()
		}
	}

	return pokemons, nil
}

// Checks an ID for every Pokemon given and compares it with those given in JSON file
func select_pokemon_by_id(id int, all_pokemons Pokemons) (Pokemon, error) {
	for j := range all_pokemons.Pokemons {
		if id == all_pokemons.Pokemons[j].ID {
			return all_pokemons.Pokemons[j], nil
		}
	}

	error_text := fmt.Sprintf("pokemon id not found; id given %d could not be found\n", id)
	err := errors.New(error_text)
	return all_pokemons.Pokemons[0], err
}

// Checks if the pokemon group being searched by the user exists or not
func is_in_pokegroups(word string) bool {

	//Valid Pokemon groups to search for in config JSON file
	var poke_groups = []string{"starter", "favorite", "legendary", "pseudolegendary", "none"}

	for _, group := range poke_groups {
		if strings.ToLower(word) == group {
			return true
		}
	}

	return false
}

// Check if ID given has been already taken by the program in previous loops
func check_if_ID_is_already_selected(id int, id_list []int) bool {
	for _, item := range id_list {
		if item == id {
			return true
		}
	}
	return false
}

// Check if pokemon type is already selected in previous loops
func is_string_already_selected(poketype string, poke_types []string) bool {
	for _, item := range poke_types {
		if strings.ToLower(poketype) == item {
			return true
		}
	}
	return false
}

// Receives a gen of Pokemon as input. Pokemon generation decides the range to iterate the random number.
// Returns an int which is the ID of the Pokemon selected
func choose_random_pokemon_by_gen(gen int) int {

	var min, max int

	rand.Seed(time.Now().UnixNano())

	switch gen {

	//Selects Pokemon from ANY generation
	case 0:
		min = 1
		max = 905

	//Selects Pokemon from 1st generation
	case 1:
		min = 1
		max = 151

	//Selects Pokemon from 2nd generation
	case 2:
		min = 152
		max = 251

	//Selects Pokemon from 3rd generation
	case 3:
		min = 252
		max = 386

	//Selects Pokemon from 4th generation
	case 4:
		min = 387
		max = 493

	//Selects Pokemon from 5th generation
	case 5:
		min = 494
		max = 649

	//Selects Pokemon from 6th generation
	case 6:
		min = 650
		max = 721

	//Selects Pokemon from 7th generation
	case 7:
		min = 722
		max = 809

	//Selects Pokemon from 8th generation
	case 8:
		min = 810
		max = 905

	//If the user has not provided a valid generation then use gen 0 as default, i.e., any value from any gen
	default:
		fmt.Println("Selected gen is not valid. Using any generation (random) instead")
		min = 1
		max = 905
	}

	return rand.Intn(max-min+1) + min //Returns an integer number between min and max
}

// Function to read big size JSON files. If some JSON files have too many variables/data in Go
// code could break. This fixes that
func read_json_big_file(pokemons *Pokemons, filename string) error {

	jsonFile, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	dec := json.NewDecoder(jsonFile)

	// read open bracket 1
	_, err_dec := dec.Token()
	if err_dec != nil {
		return err_dec
	}

	// read open bracket 2
	_, err_dec2 := dec.Token()
	if err_dec2 != nil {
		return err_dec2
	}

	// read open bracket 3
	_, err_dec3 := dec.Token()
	if err_dec3 != nil {
		return err_dec3
	}

	// while the array contains values
	for dec.More() {
		var p Pokemon
		// decode an array value
		err := dec.Decode(&p)
		if err != nil {
			return err
		}

		*&pokemons.Pokemons = append(*&pokemons.Pokemons, p)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}

	return nil
}

// Reads json file 'filename' and updates value from 'pokemons' variable with those given by json file itself
// Only useful for small JSON files; otherwise read read_json_big_file function for large JSON files
func read_json_pokemon(pokemons *Pokemons, filename string) error {

	jsonFile, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	//Read opened json file as a byte array
	byteJson, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	//Write data into 'pokemons' variable using pointers
	unmarsh_err := json.Unmarshal(byteJson, &pokemons)
	if unmarsh_err != nil {
		return unmarsh_err
	}

	return nil
}

// Reads configuration JSON file
func read_json_default_config(config *Config, filename_config string) error {

	jsonFile, err := os.Open(filename_config)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	//Read opened json file as a byte array
	byteJson, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	//Write data into 'pokemons' variable using pointers
	unmarsh_err := json.Unmarshal(byteJson, &config)
	if unmarsh_err != nil {
		return unmarsh_err
	}

	return nil
}

// Reads custom Pokemon selected in a custom file located in '/custom' folder
func read_json_custom_pokemon(custom_pokemons *Custom_pokemons, filename string) error {

	jsonFile, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	//Read opened json file as a byte array
	byteJson, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	//Write data into 'pokemons' variable using pointers
	unmarsh_err := json.Unmarshal(byteJson, &custom_pokemons)
	if unmarsh_err != nil {
		return unmarsh_err
	}

	return nil
}

// Gets the path from where the script is running. This is useful to get some data directories into the future for the program
func Get_exec_path() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

// Checks if a file exists or not. If it does exist it returns 'true'; if it does not exist returns a false.
func check_if_file_or_directory_exists(path_name string) (bool, error) {
	_, err := os.Stat(path_name)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Checks if the directory where the custom file is located exists. If not, it creates a folder called 'custom'.
// Then the program exits with code 1 (error)
func check_custom_files(path_to_custom_file string) error {

	custom_exists, err := check_if_file_or_directory_exists(path_to_custom_file)
	if err != nil {
		return err
	}

	if !custom_exists {
		err := os.Mkdir(path_to_custom_file, os.ModePerm) //directory with 0777 permission
		if err != nil {
			return err
		}
		error_text := fmt.Sprintf("no folder with location %v was found; new folder with that name created. please retry\n", path_to_custom_file)
		return errors.New(error_text)
	}
	return nil
}

// Transforms Pokemons from Custom file (which are Custom_pokemon variable type) to general data Pokemon type
func pass_custom_to_normal_Pokemon(custom_pokemons Custom_pokemons, all_pokemons Pokemons) (Pokemons, []string) {
	var (
		pokemons       Pokemons
		pokemon_styles []string
	)

	for i := range custom_pokemons.Custom_pokemons {
		for j := range all_pokemons.Pokemons {
			if strings.ToLower(custom_pokemons.Custom_pokemons[i].Name) == strings.ToLower(all_pokemons.Pokemons[j].Name) {
				new_pokemon := all_pokemons.Pokemons[j]

				new_type, valid := decide_color_custom(custom_pokemons.Custom_pokemons[i].Color)

				if valid {
					new_pokemon.Type = new_type
				}
				pokemons.Pokemons = append(pokemons.Pokemons, new_pokemon)

				continue
			}
		}
		pokemon_styles = append(pokemon_styles, custom_pokemons.Custom_pokemons[i].Style)
	}
	return pokemons, pokemon_styles
}

// Every Pokemon type has asigned a color. Therefore, for coding purposes, we change the type of the Pokemon
// E.g., if a Wailord is 'fire' type (even when it is obviously not) the code will interpret fire as red color to print
func decide_color_custom(which_color string) (string, bool) {
	var poketype string
	var valid bool

	switch strings.ToLower(which_color) {

	case "red":
		poketype, valid = "fire", true

	case "green":
		poketype, valid = "bug", true

	case "brown":
		poketype, valid = "ground", true

	case "blue":
		poketype, valid = "water", true

	case "purple":
		poketype, valid = "poison", true

	case "light_gray", "lightgray", "light gray":
		poketype, valid = "steel", true

	case "dark_gray", "darkgray", "dark gray":
		poketype, valid = "dark", true

	case "light_red", "lightred", "light red":
		poketype, valid = "fairy", true

	case "light_green", "lightgreen", "light green":
		poketype, valid = "grass", true

	case "yellow":
		poketype, valid = "electric", true

	case "light_blue", "lightblue", "light blue":
		poketype, valid = "flying", true

	case "light_purple", "lightpurple", "light purple":
		poketype, valid = "psychic", true

	case "light_cyan", "lightcyan", "light cyan":
		poketype, valid = "ice", true

	case "cyan":
		poketype, valid = "ice2", true

	default:
		poketype, valid = "", false
	}

	return poketype, valid
}

// Return a list based on styles selected on custom selection file
func create_Pokemon_styles(custom_pokemons Custom_pokemons) []string {
	var pokemon_styles []string

	for j := range custom_pokemons.Custom_pokemons {
		pokemon_styles = append(pokemon_styles, custom_pokemons.Custom_pokemons[j].Style)
	}
	return pokemon_styles
}

// Gets the color for every Pokemon based on their type
func Get_Pokemon_Colors(config Config, pokemons Pokemons) ([]string, error) {

	//Simplest case: user is using default config and provides color in config file itself
	if strings.ToLower(config.ChoosePokemon) == "default" && strings.ToLower(config.Color) != "default" && strings.ToLower(config.Color) != "random" && check_if_color_valid(config.Color) {
		return []string{strings.ToLower(config.Color), strings.ToLower(config.Color), strings.ToLower(config.Color)}, nil
	}

	//If the user uses default config file and default config colors then colors will be based on Pokemon type
	if strings.ToLower(config.ChoosePokemon) == "default" && strings.ToLower(config.Color) == "default" {
		var color_list []string
		err := colors_based_on_type(&color_list, pokemons, config.WhiteTerminal)
		if err != nil {
			return color_list, err
		}
		return color_list, nil
	}

	//If the user wants random colors choose them, considering if they can be repeated or not
	if strings.ToLower(config.ChoosePokemon) == "default" && strings.ToLower(config.Color) == "random" {
		color_list, err := randomColor(config.RepeatColor, config.WhiteTerminal)
		if err != nil {
			return color_list, err
		}
		return color_list, nil
	}

	return []string{}, nil
}

// Checks if the color given is a valid color for the code
func check_if_color_valid(possible_color string) bool {

	check_misspellings(&possible_color)

	for _, color := range valid_colors() {
		if possible_color == color {
			return true
		}
	}
	return false
}

// Check for possible words for the colors the user could bring to config or custom file
func check_misspellings(word *string) {

	switch strings.ToLower(*word) {
	case "lightblue", "light blue", "lb", "light_blue", "lblue", "l blue", "l_blue":
		*word = "light_blue"

	case "red", "redd", "r":
		*word = "red"

	case "green", "g", "greeen", "gren":
		*word = "green"

	case "brown", "br", "brow", "bron":
		*word = "brown"

	case "blue", "b", "blu":
		*word = "blue"

	case "purple", "purpl", "purp", "prple":
		*word = "purple"

	case "cyan", "cya", "cyann", "cyyan":
		*word = "cyan"

	case "light_gray", "lightgray", "light gray", "lgray", "l gray", "l_gray":
		*word = "light_gray"

	case "dark_gray", "dark gray", "darkgray", "d gray", "dgray", "d_gray":
		*word = "dark_gray"

	case "light_red", "light red", "lightred", "lr", "l red", "lred", "l_red":
		*word = "light_red"

	case "light_green", "lightgreen", "light green", "lg", "lgreen", "l green", "l_green":
		*word = "light_green"

	case "yellow", "y", "yelow", "yello", "yel":
		*word = "yellow"

	case "light_purple", "light purple", "lightpurple", "lpur", "l purple", "lpurple", "l_purple":
		*word = "light_purple"

	case "white", "w", "wite", "wt":
		*word = "white"

	default:
		*word = "default"
	}
}

// Select a color based on Pokemon type
func colors_based_on_type(color_list *[]string, pokemons Pokemons, isBackgroundConsoleWhite bool) error {
	for i := range pokemons.Pokemons {
		switch strings.ToLower(pokemons.Pokemons[i].Type) {

		case "fire":
			*color_list = append(*color_list, "red")

		case "bug":
			*color_list = append(*color_list, "green")

		case "electric":
			*color_list = append(*color_list, "yellow")

		case "water":
			*color_list = append(*color_list, "blue")

		case "poison":
			*color_list = append(*color_list, "purple")

		case "ghost":
			*color_list = append(*color_list, "purple")

		case "steel":
			*color_list = append(*color_list, "light_gray")

		case "dark":
			*color_list = append(*color_list, "dark_gray")

		case "psychic":
			*color_list = append(*color_list, "light_purple")

		case "grass":
			*color_list = append(*color_list, "light_green")

		case "dragon":
			*color_list = append(*color_list, "purple")

		case "flying":
			*color_list = append(*color_list, "light_blue")

		case "ice":
			*color_list = append(*color_list, "light_cyan")

		case "ground":
			*color_list = append(*color_list, "brown")

		case "ice2":
			*color_list = append(*color_list, "cyan")

		case "fairy":
			*color_list = append(*color_list, "light_red")

		case "fighting":
			*color_list = append(*color_list, "brown")

		case "rock":
			*color_list = append(*color_list, "light_gray")

		case "normal":
			if isBackgroundConsoleWhite {
				*color_list = append(*color_list, "black")
			} else {
				*color_list = append(*color_list, "white")
			}

		default:
			text_error := fmt.Sprintf("Pokemon type not recognized: %q is not a valid type\n", strings.ToLower(pokemons.Pokemons[i].Type))
			return errors.New(text_error)
		}
	}

	if len(*color_list) != 3 {
		return errors.New("color list has a length not supported")
	}

	return nil
}

// Select random colors from the available ones, with the choice of repeating them (check config file) and checking if
// the White Terminal option is ON/OFF
func randomColor(repeatColor, isBackgroundConsoleWhite bool) ([]string, error) {
	var (
		colors_chosen []string
		color         string
		counter       int
	)

	for {

		if counter == 3 {
			break
		}

		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(valid_colors()))

		color = valid_colors()[index]

		if color == "white" && isBackgroundConsoleWhite {
			color = "black"
		}

		if color == "black" && !isBackgroundConsoleWhite {
			color = "white"
		}

		if !repeatColor {
			if !is_string_already_selected(color, colors_chosen) {
				colors_chosen = append(colors_chosen, color)
				counter += 1
			} else {
				continue
			}
		} else {
			colors_chosen = append(colors_chosen, color)
			counter += 1
		}
	}

	if len(colors_chosen) != 3 {
		return colors_chosen, errors.New("random color list has a not supported length")
	}

	return colors_chosen, nil
}
