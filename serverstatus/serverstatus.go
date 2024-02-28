package serverstatus

import (
	"app/input"
	"app/interactions"
	"app/screen"
	"app/screen/colors"
	"strings"
)

func HandleOn() {
	printDefaultSelectionScreen(true)
	// get user input and check if the choice is valid
	choice := input.GetChar()
	if !strings.Contains("2", choice) {
		screen.Fatalln("Error: invalid choice")
	}

	// at this point, only "2" must have been
	// pressed, so we can avoid checking
	interactions.ViewLog()
}

func HandleOff() {
	printDefaultSelectionScreen(false)
	// get user input and check if the choice is valid
	choice := input.GetChar()
	if !strings.Contains("123", choice) {
		screen.Fatalln("Error: invalid choice")
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

func printDefaultSelectionScreen(isServerOn bool) {
	screen.Clear()

	if isServerOn {
		// The server is currently: ON
		// 2. View the log until the last upload

		screen.Println("The server is currently %v", colors.GreenBold("ON"))
		screen.Println("%v View the log until the last upload", colors.Bold("2."))
	} else {
		// The server is currently: OFF
		// 1. Start the server
		// 2. View the full log
		// 3. (DANGEROUS) Force upload your version of the server as the latest

		screen.Println("The server is currently %v", colors.RedBold("OFF"))
		screen.Println("%v Start the server", colors.Bold("1."))
		screen.Println("%v View the full log", colors.Bold("2."))
		screen.Println("%v (DANGEROUS) Force upload your version of the server as the latest", colors.Bold("3."))
	}
	screen.Println("")
}
