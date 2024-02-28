package server

import (
	"app/files"
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

func Download() {
	screen.Clear()
	fmt.Println("Checking local files against remote ones...")

	// first we compress the current server; this will be
	// useful to calculate the md5 hash
	if err := files.CompressGZip(ServerFolder, ServerFolderPacked); err != nil {
		fmt.Printf("Error while compressing local server: %v\n", err)
		os.Exit(1)
	}

	filesToCheck := []string{
		ServerFolderPacked,
		LogFile,
		CommandFile,
	}
	var (
		doHashesMatch = true
		hashesError   error
	)
	for _, fileToCheck := range filesToCheck {
		localHash, err := files.CalculateFileMD5(fileToCheck)
		if err != nil {
			fmt.Printf("Error while getting information about remote server: %v\n", err)
			os.Exit(1)
		}
		remoteHash, err := gdrive.GetMD5Checksum(RemoteFolder, fileToCheck)

		if remoteHash != localHash {
			doHashesMatch = false
		}
		if err != nil {
			hashesError = err
		}
	}

	// we compare the hashes; if they are different,
	// we take for granted that the remote file is newer;
	// if err != nil, it means that the file could not be
	// found, thus we cannot download the server file and
	// must use the local one
	if !doHashesMatch && hashesError == nil {
		// we first remove the previously created tarball
		if err := os.Remove(ServerFolderPacked); err != nil {
			fmt.Printf("Error while removing the compressed local server: %v\n", err)
			os.Exit(1)
		}

		// we download the new one from google drive
		screen.Clear()
		fmt.Println("Downloading latest compressed server from Google Drive...")
		// TODO: change the done into current or vice-versa
		if err := gdrive.DownloadFile(RemoteFolder, ServerFolderPacked, ServerFolderPacked, func(total, done int64) {
			// TODO
			fmt.Printf("Progress: %v%%\n", float32(done)/float32(total)*100)
		}); err != nil {
			fmt.Printf("Error while downloading remote server: %v\n", err)
			os.Exit(1)
		}

		// then we download the latest logfile from
		// google drive as well
		screen.Clear()
		fmt.Println("Downloading latest logfile from Google Drive...")
		if err := gdrive.DownloadFile(RemoteFolder, LogFile, LogFile, func(total, done int64) {
			// TODO
			fmt.Printf("Progress: %v%%\n", float32(done)/float32(total)*100)
		}); err != nil {
			fmt.Printf("Error while downloading logfile: %v\n", err)
			os.Exit(1)
		}

		// we do NOT want to download the command file,
		// as it would most certainly lead to abuse
	}

	// now we delete the previous server folder
	fmt.Println("Deleting previous server files...")
	if err := os.RemoveAll(ServerFolder); err != nil {
		fmt.Printf("Error while removing old server: %v\n", err)
		os.Exit(1)
	}

	// eventually, we decompress the newly downloaded server
	screen.Clear()
	fmt.Println("Decompressing server folder...")
	if err := files.DecompressGZip(ServerFolderPacked, ServerFolder); err != nil {
		fmt.Printf("Error while decompressing server: %v\n", err)
		os.Exit(1)
	}

	// before telling google drive that the server is on,
	// we have to make sure that the folder exists
	fmt.Println("Ensuring files and folders are in place...")
	if err := createRemoteFolderIfNotExists(); err != nil {
		fmt.Printf("Error while creating folder: %v\n", err)
		os.Exit(1)
	}

	// now we write the lockfile with content ON
	screen.Clear()
	fmt.Println("Telling Google Drive that the server is now ON...")
	if err := gdrive.WriteFileContent(RemoteFolder, LockFile, []byte("ON")); err != nil {
		fmt.Printf("Error while uploading server information to remote server: %v\n", err)
		os.Exit(1)
	}

	// and lastly we delete the packed server
	if err := os.Remove(ServerFolderPacked); err != nil {
		fmt.Printf("Error while removing the compressed local server: %v\n", err)
		os.Exit(1)
	}
}

func Upload() {
	// we first want to compress the local server folder
	screen.Clear()
	fmt.Println("Compressing server folder...")
	if err := files.CompressGZip(ServerFolder, ServerFolderPacked); err != nil {
		fmt.Printf("Error while compressing local server: %v\n", err)
		os.Exit(1)
	}

	// after that, we want to check if the server folder
	// exists in the google drive space; if it does,
	// then nothing else needs to be done; if it does not,
	// then we need to create it
	fmt.Println("Ensuring files and folders are in place...")
	if err := createRemoteFolderIfNotExists(); err != nil {
		fmt.Printf("Error while creating folder: %v\n", err)
		os.Exit(1)
	}

	// now it's time to upload the previously compressed
	// folder
	screen.Clear()
	fmt.Println("Uploading compressed server...")
	if err := gdrive.UploadFile(RemoteFolder, ServerFolderPacked, func(total, done int64) {
		// TODO
		fmt.Printf("Progress: %v%%\n", float32(done)/float32(total)*100)
	}); err != nil {
		fmt.Printf("Error while uploading local server to Google Drive: %v\n", err)
		os.Exit(1)
	}

	// then, of course, the logfile
	screen.Clear()
	fmt.Println("Uploading logfile...")
	if err := gdrive.UploadFile(RemoteFolder, LogFile, func(total, done int64) {
		fmt.Printf("Progress: %v%%\n", float32(done)/float32(total)*100)
	}); err != nil {
		fmt.Printf("Error while uploading local logfile to Google Drive: %v\n", err)
		os.Exit(1)
	}

	// then the lockfile, setting it to OFF
	screen.Clear()
	fmt.Println("Telling Google Drive that the server is now OFF...")
	if err := gdrive.WriteFileContent(RemoteFolder, LockFile, []byte("OFF")); err != nil {
		fmt.Printf("Error while uploading server information to remote server: %v\n", err)
		os.Exit(1)
	}

	// lastly, we delete the previously compressed server
	// folder, to clean up space
	if err := os.Remove(ServerFolderPacked); err != nil {
		fmt.Printf("Error while removing the compressed local server: %v\n", err)
		os.Exit(1)
	}
}

func BackupExisting() {
	// we do not care if it works or not, since we really
	// don't know if the file exists in the first place
	fmt.Println("Creating a backup of the server...")
	_ = gdrive.RenameFile(RemoteFolder, ServerFolderPacked, fmt.Sprintf("server-backup-%v.tar.gz", time.Now().Unix()))
}

func createRemoteFolderIfNotExists() error {
	if _, err := gdrive.GetFolderByName("", RemoteFolder); err != nil {
		// if err != nil, likely the folder does not exist
		// there could have been any other error but we check
		// it later when creating the folder

		if err = gdrive.CreateFolder(RemoteFolder); err != nil {
			return err
		}
	}
	return nil
}
