package interactions

import (
	"app/files"
	"app/screen"
	"app/server"
)

func StartServer() {
	// to start the server, we first need to download
	// the necessary files from google drive
	server.Download()

	// once done that, we can start the server; it does,
	// however, cause issues with the interactive mode,
	// so we need to stop it while the server is running
	// so that the server command can display things
	// without any problems
	screen.StopInteractive()
	server.Start()
	screen.StartInteractive()

	// once the server has stopped, we need one more step
	// before uploading: we want to make sure to not lose
	// any data, so we backup the existing server files
	server.BackupExisting()

	// finally, we upload the server to google drive
	server.Upload()
}

func ViewLog() {
	server.DownloadRemoteLogFile()
	files.OpenInDefaultApplication(server.DownloadedRemoteLogFile)
}

func ForceUploadServer() {
	server.Upload()
}
