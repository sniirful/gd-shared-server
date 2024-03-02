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
	if !strings.Contains("24", choice) {
		screen.Fatalln("Quitting...")
	}

	switch choice {
	case "2":
		interactions.ViewLog()
	case "4":
		interactions.ForceServerOff()
	}
}

func HandleOff() {
	printDefaultSelectionScreen(false)
	// get user input and check if the choice is valid
	choice := input.GetChar()
	if !strings.Contains("123", choice) {
		screen.Fatalln("Quitting...")
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
		// 4. Force the server to be considered OFF

		screen.Println("The server is currently %v", colors.GreenBold("ON"))
		screen.Println("%v View the log until the last upload", colors.Bold("2."))
		screen.Println("%v (DANGEROUS) Force the server to be considered OFF", colors.Bold("4."))
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
