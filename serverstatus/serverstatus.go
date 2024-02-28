package serverstatus

import (
	"app/input"
	"app/interactions"
	"app/screen"
	"fmt"
	"os"
	"strings"
)

func HandleOff() {
	screen.PrintDefaultSelectionScreen(false)
	// get user input and check if the choice is valid
	choice := input.GetChar()
	if !strings.Contains("123", choice) {
		fmt.Println("Error: invalid choice")
		os.Exit(1)
	}

	switch choice {
	case "1":
		interactions.StartServer()
	case "2":
		interactions.ViewLog()
	case "3":
		interactions.ForceUploadServer()
	}
}

func HandleOn() {
	screen.PrintDefaultSelectionScreen(true)
	// get user input and check if the choice is valid
	choice := input.GetChar()
	if !strings.Contains("2", choice) {
		fmt.Println("Error: invalid choice")
		os.Exit(1)
	}

	// at this point, only "2" must have been
	// pressed, so we can avoid checking
	interactions.ViewLog()
}
