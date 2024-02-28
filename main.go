package main

import (
	"app/gdrive"
	"app/screen"
	"app/server"
	"app/serverstatus"
	"fmt"
)

func main() {
	screen.Clear()
	fmt.Println("Connecting to Google Drive...")

	// check if the app can connect to google drive
	if _, err := gdrive.ListAllFiles(); err != nil {
		fmt.Printf("Error while connecting to Google Drive: %v\n", err)
		return
	}

	// check if the server is up:
	// it is off if the lockfile does not exist (err != nil)
	// or if it exists and contains anything other than "ON"
	// inside of it
	if bytes, err := gdrive.GetFileContent(server.RemoteFolder, server.LockFile); err != nil || string(bytes) != "ON" {
		serverstatus.HandleOff()
	} else {
		serverstatus.HandleOn()
	}
}
