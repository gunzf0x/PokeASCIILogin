package print_color

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// Styles for text to print
const (
	normal  string = ""
	blink   string = "5;"
	reverse string = "7;"
)

// Type that will be used to get terminal size
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Sub-main function. Prints ascii art in colors based on the user's decision given in config file
func Print_in_colors(widths_asciis []int, print_style []string, ascii_art_to_print string,
	colors_selected []string) string {

	var (
		ascii_in_color string
		centering      string
	)

	lines_to_print := strings.Split(ascii_art_to_print, "\n")

	maxWidth := get_max_width(lines_to_print)

	colsTerminal, err_cols := getTerminalCols()
	if err_cols != nil {
		centering = ""
	} else {
		// Center the text adding spaces
		for i := 0; i < (int(colsTerminal)-maxWidth)/2; i++ {
			centering += " "
		}
	}

	for _, line := range lines_to_print {
		if len(line) == 0 {
			continue
		}
		sub_line1 := line[:widths_asciis[0]]
		sub_line2 := line[widths_asciis[0]:(widths_asciis[0] + widths_asciis[1])]
		sub_line3 := line[(widths_asciis[0] + widths_asciis[1]):(widths_asciis[0] + widths_asciis[1] + widths_asciis[2])]

		ascii_in_color += centering
		ascii_in_color += add_color_string(sub_line1, colors_selected[0], print_style[0]) + add_color_string(sub_line2, colors_selected[1], print_style[1]) + add_color_string(sub_line3, colors_selected[2], print_style[2])
		ascii_in_color += "\n"
	}

	return ascii_in_color
}

func getTerminalCols() (uint, error) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return 0, errno
	}
	return uint(ws.Col), nil
}

func get_max_width(lines_list []string) int {

	var max_width int

	for _, line := range lines_list {
		if len(line) > max_width {
			max_width = len(line)
		}
	}
	return max_width
}

// Adds a color and a style to a string
func add_color_string(text_to_add, color, style string) string {

	var (
		add_text    string
		reset_color string
	)

	reset_color = string("\033[0m")

	//Select style for ascii art
	switch strings.ToLower(style) {
	case "blink":
		style = blink

	case "reverse":
		style = reverse

	case "none", "", "normal", "default":
		style = normal

	default:
		style = normal
	}
	//Define colors that will be used in this program. ANSI colors
	switch strings.ToLower(color) {

	case "red":
		add_text = string(string("\033["+style+"31m") + text_to_add + reset_color)

	case "green":
		add_text = string("\033["+style+"32m") + text_to_add + reset_color

	case "brown":
		add_text = string("\033["+style+"33m") + text_to_add + reset_color

	case "blue":
		add_text = string("\033["+style+"34m") + text_to_add + reset_color

	case "purple":
		add_text = string("\033["+style+"35m") + text_to_add + reset_color

	case "cyan":
		add_text = string("\033["+style+"36m") + text_to_add + reset_color

	case "light_gray":
		add_text = string("\033["+style+"37m") + text_to_add + reset_color

	case "dark_gray":
		add_text = string("\033["+style+"1;30m") + text_to_add + reset_color

	case "light_red":
		add_text = string("\033["+style+"1;31m") + text_to_add + reset_color

	case "light_green":
		add_text = string("\033["+style+"1;32m") + text_to_add + reset_color

	case "yellow":
		add_text = string("\033["+style+"1;33m") + text_to_add + reset_color

	case "light_blue":
		add_text = string("\033["+style+"1;34m") + text_to_add + reset_color

	case "light_purple":
		add_text = string("\033["+style+"1;35m") + text_to_add + reset_color

	case "light_cyan":
		add_text = string("\033["+style+"1;36m") + text_to_add + reset_color

	case "white":
		add_text = string("\033["+style+"1;37m") + text_to_add + reset_color

	case "black":
		add_text = string("\033["+style+"30m") + text_to_add + reset_color

	default:
		add_text = text_to_add
	}

	return add_text
}

// Select styles for text such as text blinking, reverse background text with color text, none; among others
func select_text_style(style *string, config_style string) {

	config_style = strings.ToLower(config_style)

	if config_style == "none" || config_style == "normal" || config_style == "" || config_style == "default" {
		*style = normal
		return
	}

	if config_style == "blink" {
		*style = blink
		return
	}

	if config_style == "reverse" {
		*style = reverse
		return
	}

	fmt.Printf("Warning! %q is not a valid format. Using default style value instead (\"none\" default-style)\n", config_style)
	*style = normal
}
