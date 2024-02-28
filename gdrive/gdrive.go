package gdrive

import (
	"app/files"
	"app/gdrive/gdriveservice"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"google.golang.org/api/drive/v3"
)

const (
	MimeTypeFolder = "application/vnd.google-apps.folder"
)

func ListAllFiles() ([]*drive.File, error) {
	return listFiles("trashed = false")
}

func ListFilesInFolder(foldername string) ([]*drive.File, error) {
	folder, err := GetFolderByName("", foldername)
	if err != nil {
		return nil, err
	}

	return listFiles(fmt.Sprintf("trashed = false and '%v' in parents", folder.Id))
}

func RenameFile(parentFolder, oldFileName, newFileName string) error {
	service, err := gdriveservice.GetService()
	if err != nil {
		return err
	}
	// we ignore the error because we later ignore
	// the folder if it is equal to nil
	folder, _ := GetFolderByName("", parentFolder)

	driveFile, err := GetFileByName(parentFolder, oldFileName)
	if err != nil {
		return err
	}
	newFile := &drive.File{
		Name: newFileName,
	}
	// if the folder actually exists, add the file
	// to that folder, but only if we are creating
	// a new file, that is if err != nil
	if folder != nil && err != nil {
		newFile.Parents = []string{folder.Id}
	}

	_, err = service.Files.Update(driveFile.Id, newFile).Do()
	return err
}

func GetFileByName(parentFolder, filename string) (*drive.File, error) {
	file, err := getFileByFunction(parentFolder, func(f *drive.File) bool {
		return f.Name == filename && f.MimeType != MimeTypeFolder
	})
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("file not found: %v", filename)
	}

	return file, nil
}

func GetFolderByName(parentFolder, filename string) (*drive.File, error) {
	folder, err := getFileByFunction(parentFolder, func(f *drive.File) bool {
		return f.Name == filename && f.MimeType == MimeTypeFolder
	})
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, fmt.Errorf("folder not found: %v", filename)
	}

	return folder, nil
}

func CreateFolder(foldername string) error {
	service, err := gdriveservice.GetService()
	if err != nil {
		return err
	}

	folder := &drive.File{
		Name:     foldername,
		MimeType: MimeTypeFolder,
	}
	_, err = service.Files.Create(folder).Do()

	return err
}

func GetFileContent(parentFolder, filename string) ([]byte, error) {
	res, _, err := getFileDownloadResponse(parentFolder, filename)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func WriteFileContent(parentFolder, filename string, content []byte) error {
	reader := bytes.NewReader(content)
	return uploadFileReader(parentFolder, filename, reader, func(_, _ int64) {})
}

func GetMD5Checksum(parentFolder, filename string) (string, error) {
	// we get the file and check exclusively if it is a file
	file, err := GetFileByName(parentFolder, filename)
	if err != nil {
		return "", err
	}
	if file == nil {
		return "", fmt.Errorf("file not found: %v", filename)
	}

	return file.Md5Checksum, nil
}

func DownloadFile(parentFolder, remoteFileName, localFileName string, progressFunction func(total, current int64)) error {
	res, f, err := getFileDownloadResponse(parentFolder, remoteFileName)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	outFile, err := os.Create(localFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err = files.Copy(outFile, res.Body, func(done int64) {
		progressFunction(f.Size, done)
	}); err != nil {
		return err
	}

	return nil
}

func UploadFile(parentFolder, filename string, progressFunction func(total, current int64)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return uploadFileReader(parentFolder, filename, file, func(current, total int64) {
		progressFunction(total, current)
	})
}

func listFiles(query string) ([]*drive.File, error) {
	service, err := gdriveservice.GetService()
	if err != nil {
		return nil, err
	}

	// we get the files':
	// - id
	// - name
	// - size
	// - md5Checksum
	// - mimeType
	// and all other parameters are ignored, to save
	// bandwidth
	fileList, err := service.Files.List().
		Q(query).
		Fields("files(id, name, size, md5Checksum, mimeType)").
		Do()
	if err != nil {
		return nil, err
	}

	return fileList.Files, nil
}

func uploadFileReader(parentFolder, filename string, reader io.Reader, progressFunction func(total, current int64)) error {
	service, err := gdriveservice.GetService()
	if err != nil {
		return err
	}
	// we ignore the error because we later ignore
	// the folder if it is equal to nil
	folder, _ := GetFolderByName("", parentFolder)

	// we first attempt to get the already existing file;
	// if we get an error, it must mean that it does not
	// exist, thus we need to create it
	driveFile, err := GetFileByName(parentFolder, filename)
	newFile := &drive.File{
		Name: filename,
	}
	// if the folder actually exists, add the file
	// to that folder, but only if we are creating
	// a new file, that is if err != nil
	if folder != nil && err != nil {
		newFile.Parents = []string{folder.Id}
	}

	// if there was an error before, we have to create a
	// new file; if there was not, then we only have to
	// update the existing one
	if err != nil {
		_, err = service.Files.Create(newFile).Media(reader).ProgressUpdater(func(current, total int64) {
			progressFunction(total, current)
		}).Do()
	} else {
		_, err = service.Files.Update(driveFile.Id, newFile).Media(reader).ProgressUpdater(func(current, total int64) {
			progressFunction(total, current)
		}).Do()
	}
	return err
}

func getFileDownloadResponse(parentFolder, filename string) (*http.Response, *drive.File, error) {
	// we get the file and check exclusively if it is a file
	file, err := GetFileByName(parentFolder, filename)
	if err != nil {
		return nil, nil, err
	}

	service, err := gdriveservice.GetService()
	if err != nil {
		return nil, nil, err
	}

	res, err := service.Files.Get(file.Id).Download()
	return res, file, err
}

func getFileByFunction(parentFolder string, testFunction func(*drive.File) bool) (*drive.File, error) {
	var (
		files []*drive.File
		err   error
	)
	// if the parentFolder does not exist (""), then
	// get all files regardless; if it does, get all
	// files in that folder
	if parentFolder == "" {
		files, err = ListAllFiles()
	} else {
		files, err = ListFilesInFolder(parentFolder)
	}
	if err != nil {
		return nil, err
	}

	// for each file, run the test function and if it
	// passes, return the file
	for _, f := range files {
		if testFunction(f) {
			return f, nil
		}
	}
	return nil, nil
}
