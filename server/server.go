package server

import (
	"app/commands"
	"app/files"
	"app/files/fileflags"
	"app/files/filemodes"
	"app/gdrive"
	"app/screen"
	"fmt"
	"os"
	"time"
)

const (
	ServerFolder       = "server"
	ServerFolderPacked = "server.tar.gz"
	RemoteFolder       = "GD-Server"

	LogFile                 = "logfile"
	DownloadedRemoteLogFile = "logfile.remote"
	LockFile                = "lockfile"
	CommandFile             = "start.command"
)

func IsOn() bool {
	// check if the server is up:
	// it is off if the lockfile does not exist (err != nil)
	// or if it exists and contains anything other than "ON"
	// inside of it, so we check for that and return the
	// inverse of it
	bytes, err := gdrive.GetFileContent(RemoteFolder, LockFile)
	return !(err != nil || string(bytes) != "ON")
}

func Start() {
	screen.ClearAndPrintln("Preparing to start server...")

	screen.Println("Opening logfile in write mode...")
	logFile, err := os.OpenFile(LogFile, fileflags.APPEND, filemodes.RW_R__R__)
	if err != nil {
		screen.Fatalln("Error while opening logfile in write mode: %v", err)
	}
	defer logFile.Close()

	screen.Println("Reading command file...")
	commandBytes, err := os.ReadFile(CommandFile)
	if err != nil {
		screen.Fatalln("Error while reading command file: %v", err)
	}

	screen.Println("Starting server...")
	commands.RunWithWorkingDirAndLogFile(string(commandBytes), ServerFolder, logFile)
}

func Download() {
	screen.ClearAndPrintln("Preparing to download server from Google Drive...")

	screen.Println("Ensuring files and folders are in place...")
	if err := createRemoteFolderIfNotExists(); err != nil {
		screen.Fatalln("Error while checking on the folder in Google Drive: %v", err)
	}

	screen.Println("Checking local files against remote ones...")
	doHashesMatch, hashesError := checkFilesHashesMatch()

	// we check not only that hashes match, but also that
	// there were no errors because if there were, it means
	// that some files may have not been in place
	if !doHashesMatch && hashesError == nil {
		screen.Println("Found different files. Downloading files from Google Drive...")

		screen.Println("Downloading compressed server from Google Drive...")
		if err := gdrive.DownloadFile(RemoteFolder, ServerFolderPacked, ServerFolderPacked, screen.PrintProgress); err != nil {
			screen.Fatalln("Error while downloading compressed server from Google Drive: %v", err)
		}

		screen.Println("Downloading logfile from Google Drive...")
		if err := gdrive.DownloadFile(RemoteFolder, LogFile, LogFile, screen.PrintProgress); err != nil {
			screen.Fatalln("Error while downloading logfile from Google Drive: %v", err)
		}
	} else {
		screen.Println("Found the most recent files to be local. Skipping downloading from Google Drive...")
	}

	screen.Println("Deleting old server files...")
	if err := os.RemoveAll(ServerFolder); err != nil {
		screen.Fatalln("Error while deleting old server files: %v", err)
	}

	screen.Println("Decompressing server file...")
	if err := files.DecompressGZip(ServerFolderPacked, ServerFolder); err != nil {
		screen.Fatalln("Error while decompressing server file: %v", err)
	}

	screen.Println("Telling Google Drive that the server is now ON...")
	if err := gdrive.WriteFileContent(RemoteFolder, LockFile, []byte("ON")); err != nil {
		screen.Fatalln("Error while telling Google Drive about the new server status: %v", err)
	}

	screen.Println("Cleaning up...")
	if err := os.Remove(ServerFolderPacked); err != nil {
		screen.Fatalln("Error while cleaning up: %v", err)
	}
}

func Upload() {
	screen.ClearAndPrintln("Preparing to upload server to Google Drive...")

	screen.Println("Ensuring files and folders are in place...")
	if err := createRemoteFolderIfNotExists(); err != nil {
		screen.Fatalln("Error while checking on the folder in Google Drive: %v", err)
	}

	screen.Println("Compressing server folder...")
	if err := files.CompressGZip(ServerFolder, ServerFolderPacked); err != nil {
		screen.Fatalln("Error while compressing server folder: %v", err)
	}

	screen.Println("Uploading compressed server to Google Drive...")
	if err := gdrive.UploadFile(RemoteFolder, ServerFolderPacked, screen.PrintProgress); err != nil {
		screen.Fatalln("Error while uploading compressed server to Google Drive: %v", err)
	}

	screen.Println("Uploading logfile to Google Drive...")
	if err := gdrive.UploadFile(RemoteFolder, LogFile, screen.PrintProgress); err != nil {
		screen.Fatalln("Error while uploading logfile to Google Drive: %v", err)
	}

	screen.Println("Telling Google Drive that the server is now OFF...")
	if err := gdrive.WriteFileContent(RemoteFolder, LockFile, []byte("OFF")); err != nil {
		screen.Fatalln("Error while telling Google Drive about the new server status: %v", err)
	}

	screen.Println("Cleaning up...")
	if err := os.Remove(ServerFolderPacked); err != nil {
		screen.Fatalln("Error while cleaning up: %v", err)
	}
}

func BackupExisting() {
	// we do not care if it works or not, since we really
	// don't know if the file exists in the first place
	screen.Println("Creating a backup of the server...")
	_ = gdrive.RenameFile(RemoteFolder, ServerFolderPacked, fmt.Sprintf("server-backup-%v.tar.gz", time.Now().Unix()))
}

func DownloadRemoteLogFile() {
	screen.ClearAndPrintln("Downloading remote logfile...")
	if err := gdrive.DownloadFile(RemoteFolder, LogFile, DownloadedRemoteLogFile, screen.PrintProgress); err != nil {
		screen.Fatalln("Error while downloading remote logfile: %v", err)
	}
}

func createRemoteFolderIfNotExists() error {
	// if err != nil, likely the folder does not exist
	// there could have been any other error but we check
	// it later when creating the folder
	if _, err := gdrive.GetFolderByName("", RemoteFolder); err != nil {
		if err = gdrive.CreateFolder(RemoteFolder); err != nil {
			return err
		}
	}
	return nil
}

func checkFilesHashesMatch() (bool, error) {
	// before checking all the files, we firstly need
	// to compress the server folder, so we can check
	// its hash against the google drive one
	if err := files.CompressGZip(ServerFolder, ServerFolderPacked); err != nil {
		screen.Fatalln("Error while compressing local server: %v", err)
	}

	filesToCheck := []string{
		ServerFolderPacked,
		LogFile,
	}
	var (
		doHashesMatch = true
		hashesError   error
	)
	for _, fileToCheck := range filesToCheck {
		localHash, err := files.CalculateFileMD5(fileToCheck)
		if err != nil {
			screen.Fatalln("Error while checking file %v: %v", fileToCheck, err)
		}
		remoteHash, err := gdrive.GetMD5Checksum(RemoteFolder, fileToCheck)

		if remoteHash != localHash {
			doHashesMatch = false
		}
		if err != nil {
			hashesError = err
		}
	}

	return doHashesMatch, hashesError
}
