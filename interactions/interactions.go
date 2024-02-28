package interactions

import (
	"app/commands"
	"app/files"
	"app/gdrive"
	"app/screen"
	"app/server"
	"fmt"
	"os"
)

func StartServer() {
	// the server is first downloaded
	server.Download()

	// we first open the logfile that we will
	// give to the command; we use os.OpenFile because
	// with os.Open the file is in readonly mode
	logFile, err := os.OpenFile(server.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error while opening logfile: %v\n", err)
		os.Exit(1)
	}
	// we do not defer the logFile.Close() because
	// we need to close the file exactly after the
	// program has run

	// after that, the server is run with the
	// command taken from the file
	data, err := os.ReadFile(server.CommandFile)
	if err != nil {
		fmt.Printf("Error while reading command file: %v\n", err)
		os.Exit(1)
	}

	screen.Clear()
	fmt.Println("Starting server...")
	// we ignore SIGINT, this way when the user presses
	// Ctrl-C, only the child process will stop
	// signal.Ignore(syscall.SIGINT)
	commands.RunWithWorkingDirAndLogFile(string(data), server.ServerFolder, logFile)
	logFile.Close()
	// at this point, we reset the SIGINT behaviour
	// signal.Reset(syscall.SIGINT)

	// before uploading, we make a backup; just to be sure
	server.BackupExisting()

	// once the server has stopped, we need to
	// upload it back to google drive
	server.Upload()
}

func ViewLog() {
	screen.Clear()
	fmt.Println("Downloading remote logfile...")

	// download logfile to logfile.remote
	if err := gdrive.DownloadFile(server.RemoteFolder, server.LogFile, server.DownloadedRemoteLogFile, func(total, done int64) {
		// TODO
		fmt.Printf("Progress: %v%%\n", float32(done)/float32(total)*100)
	}); err != nil {
		fmt.Printf("Error while downloading logfile: %v\n", err)
		os.Exit(1)
	}

	// open it with default editor
	screen.Clear()
	fmt.Println("Opening downloaded logfile in default editor...")
	files.OpenInDefaultApplication(server.DownloadedRemoteLogFile)
}

func ForceUploadServer() {
	server.Upload()
}
