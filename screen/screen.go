package screen

import (
	"app/commands"
	"app/screen/colors"
	"fmt"
	"runtime"
)

func Clear() {
	var command string
	if runtime.GOOS == "windows" {
		command = "cls"
	} else {
		command = "clear"
	}

	command = ""
	commands.Run(command)
}

func PrintDefaultSelectionScreen(isServerOn bool) {
	Clear()

	if isServerOn {
		// The server is currently: ON
		// 2. View the log until the last upload

		fmt.Printf("The server is currently %v"+"\n"+
			"%v View the log until the last upload"+"\n",
			colors.GreenBold("ON"), colors.Bold("2."))
	} else {
		// The server is currently: OFF
		// 1. Start the server
		// 2. View the full log
		// 3. (DANGEROUS) Force upload your version of the server as the latest

		fmt.Printf("The server is currently %v"+"\n"+
			"%v Start the server"+"\n"+
			"%v View the full log"+"\n"+
			"%v (DANGEROUS) Force upload your version of the server as the latest"+"\n",
			colors.RedBold("OFF"), colors.Bold("1."), colors.Bold("2."), colors.Bold("3."))
	}
	fmt.Println()
}
