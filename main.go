package main

import (
	"Termify/auth"
	"Termify/helpers"
	"Termify/server"
	"bufio"
	"fmt"
	"os"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	startupText = "termify"
	exitCode    = "9"
)

func startup(text string) {
	helpers.ClearTerm()
	color.Green(banner.PrintS(text))
}

func main() {
	startup(startupText)
	auth.Authorize()
	srv := server.Create()
	server.Start(srv)
	scanner := bufio.NewScanner(os.Stdin)

	// When on the main menu, keep scanning input until the user
	// chooses to exit
	for scanner.Scan() {
		text := scanner.Text()
		if text == exitCode {
			return
		}
		fmt.Println(text)
	}
}
