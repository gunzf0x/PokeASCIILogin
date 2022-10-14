package read_files

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"pokemon_ascii/get_data"
	"strings"
)

// Concatenates the lines for 3 files that will be printed
// Search for the ascii art filenames corresponding to Pokemon names, concatenates them in one single string and return them at the same level/height
// Filename1, filename2 and filename3 will be printed from left to right, respectively
func Unify_files(path string, pokemons get_data.Pokemons) (string, []int, error) {

	var ascii_art_path string = path + "/ascii_pokemon_files/"

	//Change the name of the Pokemon so they have the same name as their ascii art files and can be readen by the program
	filename1, filename2, filename3 := clean_filenames(pokemons.Pokemons[0].Name),
		clean_filenames(pokemons.Pokemons[1].Name),
		clean_filenames(pokemons.Pokemons[2].Name)

	//Count the number of lines for each file
	n_lines_filename1, err_f1 := n_lines_file(ascii_art_path + filename1 + ".txt")
	if err_f1 != nil {
		return "", []int{}, err_f1
	}

	n_lines_filename2, err_f2 := n_lines_file(ascii_art_path + filename2 + ".txt")
	if err_f2 != nil {
		return "", []int{}, err_f2
	}

	n_lines_filename3, err_f3 := n_lines_file(ascii_art_path + filename3 + ".txt")
	if err_f3 != nil {
		return "", []int{}, err_f3
	}

	//Now check which one of them is the biggest in length, i.e., the one with more number of lines
	larger_length, err_length := max_lengths(n_lines_filename1, n_lines_filename2, n_lines_filename3)
	if err_length != nil {
		return "", []int{}, err_length
	}

	//Compute the difference of lines between each file and the file with max lines
	difference, err_diff := n_lines_of_difference(larger_length, n_lines_filename1, n_lines_filename2, n_lines_filename3)
	if err_diff != nil {
		return "", []int{}, err_diff
	}

	diff1, diff2, diff3 := difference[0], difference[1], difference[2]

	//Read ascii art from the corresponding files, files which names are Pokemon names in lowercase
	ascii_art1, err_ascii1 := read_files(ascii_art_path + filename1 + ".txt")
	if err_ascii1 != nil {
		return "", []int{}, err_ascii1
	}
	ascii_art2, err_ascii2 := read_files(ascii_art_path + filename2 + ".txt")
	if err_ascii2 != nil {
		return "", []int{}, err_ascii2
	}

	ascii_art3, err_ascii3 := read_files(ascii_art_path + filename3 + ".txt")
	if err_ascii3 != nil {
		return "", []int{}, err_ascii3
	}

	//Finally, concatenate the files, line by line, so 3 different ascii arts can be printed from left to right
	concatenate_ascii, max_lengths, err_concatenate := concatenate_lines(diff1, diff2, diff3, ascii_art1, ascii_art2, ascii_art3)
	if err_concatenate != nil {
		return "", []int{}, err_concatenate
	}

	//Finally, return ascii to be printed, the width of every ascii (so the program will know the range to print in color)
	//and if everything went ok, a nil (no-error)
	return concatenate_ascii, max_lengths, nil
}

// Reads files and gets their content. Returns a string containing all the lines of the file and an error (not nil if the filename is not valid)
func read_files(filename string) (string, error) {

	var lines string

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return line, err
		}
		line = strings.TrimRight(line, "\r\n")
		lines += line + "\n"
	}
	return lines, nil
}

// Creates one big string concatenating the first line from file 1 with the first line from file 2 and 3, and so on for the rest of the lines
func concatenate_lines(diff1, diff2, diff3 int, ascii1, ascii2, ascii3 string) (string, []int, error) {

	var (
		final_string    string
		max_length_list []int
	)

	//Add lines to all the ascii's art if they have less height than the file with max height
	ascii1, max_width_ascii1 := add_width(diff1, ascii1)
	ascii2, max_width_ascii2 := add_width(diff2, ascii2)
	ascii3, max_width_ascii3 := add_width(diff3, ascii3)

	ascii1_lines := set_all_lines_to_equal_width(ascii1, max_width_ascii1)
	ascii2_lines := set_all_lines_to_equal_width(ascii2, max_width_ascii2)
	ascii3_lines := set_all_lines_to_equal_width(ascii3, max_width_ascii3)

	//if the files have not the same height means that they will printed at different height and, therefore, not as we want
	if (len(ascii1_lines) != len(ascii2_lines)) || (len(ascii1_lines) != len(ascii3_lines)) || (len(ascii2_lines) != len(ascii3_lines)) {
		error_text := fmt.Sprintf("modified lengths of files are not the same: file 1 has %d lines, file 2 has %d lines and file 3 has %d lines\n", len(ascii1_lines), len(ascii2_lines), len(ascii3_lines))
		err := errors.New(error_text)
		return "", max_length_list, err
	}

	for j := 0; j < len(ascii1_lines); j++ {
		final_string += ascii1_lines[j] + ascii2_lines[j] + ascii3_lines[j] + "\n"
	}
	max_length_list = append(max_length_list, max_width_ascii1, max_width_ascii2, max_width_ascii3)
	return final_string, max_length_list, nil
}

// Every line of the ascii file will now have the same width which makes things easier to print ascii's art to its left
func set_all_lines_to_equal_width(ascii string, max_width int) []string {

	var list_ascii []string

	for _, line := range strings.Split(ascii, "\n") {
		for {
			if max_width > len(line) {
				for i := 0; i < (max_width - len(line)); i++ {
					line += " "
				}

			} else {
				break
			}
		}
		list_ascii = append(list_ascii, line)
	}
	return list_ascii
}

// Since the ascii arts do not have the same height, we make them to be at the same height
func add_width(diff int, ascii string) (string, int) {
	var (
		add_line  string
		max_width int
	)

	max_width = get_max_width_of_file(ascii)
	//First, add lines so all the asciis are at the same altitude
	for i := 0; i < diff; i++ {
		for j := 0; j < max_width; j++ {
			add_line += " "
		}
		add_line += "\n"
	}

	ascii = add_line + ascii

	return ascii, max_width
}

// Obtains the max width (columns) for every file
func get_max_width_of_file(ascii string) int {

	var (
		lines_ascii []string
		max         int
	)

	lines_ascii = strings.Split(ascii, "\n")

	for _, line := range lines_ascii {

		if len(line) > max {
			max = len(line)
		}

	}
	return max
}

// Returns the number of lines of a file
func n_lines_file(file_name string) (int, error) {

	f, err := os.Open(file_name)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	lines_f, err := lineCounter(f)
	if err != nil {
		return 0, err
	}
	return lines_f, err

}

// Counts lines of a file in the most eficient way possible
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// Compares numbers and returns the max of them
func max_lengths(lengths ...int) (int, error) {

	var max int

	for i, value := range lengths {
		//Since we are comparing slice lengths they cannot be negative
		if value < 0 {
			error_text := fmt.Sprintf("error in length for element %d: negative length (%d)\n", i+1, value)
			err := errors.New(error_text)
			return 0, err
		}

		if value > max {
			max = value
		}
	}

	return max, nil
}

// Returns a list with the difference between 'max' and the corresponding argument
func n_lines_of_difference(max int, lengths ...int) ([]int, error) {

	var difference []int

	for i, value := range lengths {

		if max < value {
			error_text := fmt.Sprintf("length for element %d is bigger than max value: max is %d and value is %d, which should not be possible\n", i+1, max, value)
			err := errors.New(error_text)
			return difference, err
		}

		difference = append(difference, max-value)
	}

	return difference, nil
}

// Cleans Pokemon names from spaces and parenthesis so ascii files with their names can be readen
// Parenthesis are deleted, spaces (" ") and lines ("-") are replaced by underlines (_)
// Finally string is returned in lowercase
func clean_filenames(filename_to_read string) string {

	filename_to_read = strings.ReplaceAll(filename_to_read, "(", "")
	filename_to_read = strings.ReplaceAll(filename_to_read, ")", "")
	filename_to_read = strings.ReplaceAll(filename_to_read, " ", "_")
	filename_to_read = strings.ReplaceAll(filename_to_read, "-", "_")
	filename_to_read = strings.ToLower(filename_to_read)

	return filename_to_read
}
