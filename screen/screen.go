package screen

import (
	"app/commands"
	"fmt"
	"log"
	"os"
	"runtime"

	"atomicgo.dev/cursor"
	"golang.org/x/term"
)

var (
	oldTerminalState *term.State = nil

	// we save this variable so that when we print
	// progress, if the previous log was progress it
	// moves the cursor up; if it is not, it prints
	// the progress regardless
	previousLogIsProgressLog = false
)

func StartInteractive() {
	// this "interactive" mode allows the program to
	// move the cursor freely, along with giving us the
	// raw input, that is whenever the user presses a key,
	// it's immediately passed to the program and the user
	// doesn't need to press enter
	var err error
	oldTerminalState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Unable to start interactive mode: %v\n", err)
	}
}

func StopInteractive() {
	term.Restore(int(os.Stdin.Fd()), oldTerminalState)
}

func Clear() {
	var command string
	if runtime.GOOS == "windows" {
		command = "cls"
	} else {
		command = "clear"
	}

	commands.Run(command)
}

func Println(format string, args ...any) {
	// \n\r because we assume the interactive mode, thus
	// if we go to a newline it won't actually be a newline
	// but rather it will just go down with the cursor
	fmt.Printf("%v\n\r", fmt.Sprintf(format, args...))
	previousLogIsProgressLog = false
}

func ClearAndPrintln(format string, args ...any) {
	Clear()
	Println(format, args...)
}

// TODO: decide between done and current
func PrintProgress(total, done int64) {
	if previousLogIsProgressLog {
		cursor.UpAndClear(1)
	}
	Println("Progress: %v%%", float32(done)/float32(total)*100)
	previousLogIsProgressLog = true
}

func Fatalln(format string, args ...any) {
	Println(format, args...)
	StopInteractive()
	os.Exit(1)
}
