package main

import (
	"app/gdrive"
	"app/screen"
	"app/server"
	"app/serverstatus"
	"app/signals"
)

func main() {
	screen.StartInteractive()
	defer screen.StopInteractive()

	// we want to handle Ctrl-C on our own
	onStopSignalling := signals.CaptureInterrupt(func() {
		// TODO
		screen.StopInteractive()
	})
	defer onStopSignalling()

	// by listing the files and doing nothing with them,
	// we can make sure that the connection with google
	// drive works properly
	screen.ClearAndPrintln("Connecting to Google Drive...")
	if _, err := gdrive.ListAllFiles(); err != nil {
		screen.Println("Error while connecting to Google Drive: %v", err)
		return
	}

	if server.IsOn() {
		serverstatus.HandleOn()
	} else {
		serverstatus.HandleOff()
	}
}
