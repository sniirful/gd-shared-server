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
		// there's nothing we need to do on interrupt
		// captured, we just need to make sure that the
		// program doesn't exit because the user pressed
		// Ctrl-C by mistake
	})
	defer onStopSignalling()

	// we're using an infinite loop instead of recursion
	// in order not to call the above functions
	// without a reason to do so
	for {
		// by listing the files and doing nothing with them,
		// we can make sure that the connection with google
		// drive works properly
		screen.ClearAndPrintln("Connecting to Google Drive...")
		driveUsage, driveLimit, err := gdrive.GetUsageQuota()
		if err != nil {
			screen.Println("Error while connecting to Google Drive: %v", err)
			return
		}
		// we ignore the error here because if the server hasn't
		// been set up yet, we don't want to print an error
		serverSize, _ := gdrive.GetFileSize(server.RemoteFolder, server.ServerFolderPacked)

		if server.IsOn() {
			serverstatus.HandleOn(serverSize, driveUsage, driveLimit)
		} else {
			serverstatus.HandleOff(serverSize, driveUsage, driveLimit)
		}
	}
}
