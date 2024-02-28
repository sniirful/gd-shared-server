package files

import (
	"app/commands"
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	BufferSize = 65535
)

func OpenInDefaultApplication(filename string) {
	command := fmt.Sprintf("%v %v", func() string {
		if runtime.GOOS == "windows" {
			return "start"
		} else {
			return "xdg-open"
		}
	}(), filename)
	commands.RunSilent(command)
}

func Copy(outFile *os.File, reader io.Reader, progressFunction func(current int64)) error {
	buffer := make([]byte, BufferSize)
	var progressDone int64
	for {
		read, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// we ignore the first value because it's the
		// number of bytes written, all that we gave it
		_, err = outFile.Write(buffer[:read])
		if err != nil {
			return err
		}

		progressDone += int64(read)
		progressFunction(progressDone)
	}
	return nil
}

func CalculateFileMD5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// stolen deliberately from google bard
func CompressGZip(inFolder, outFile string) error {
	// Open the output file for writing
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a new gzip compressor writer
	gw := gzip.NewWriter(out)
	defer gw.Close()

	// Create a new tar archive
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Walk through the input folder and add files to the archive
	err = filepath.Walk(inFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root folder itself
		if path == inFolder {
			return nil
		}

		// Create a header for the file
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		// Ensure path is relative to the input folder
		header.Name = strings.TrimPrefix(path, inFolder)

		// Write the header to the tar archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If it's a regular file, write its contents to the archive
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// stolen deliberately from google bard
func DecompressGZip(inFile, outFolder string) error {
	// without this, no folder will be created
	os.MkdirAll(outFolder, os.ModePerm)

	// Open the gzip file for reading
	in, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create a new gzip reader
	gr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gr.Close()

	// Create a new tar archive reader
	tr := tar.NewReader(gr)

	// Extract each file from the archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // Reached end of archive
		}
		if err != nil {
			return err
		}

		// Create the output file
		outFile := filepath.Join(outFolder, header.Name)

		// Check if it's a directory
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(outFile, 0755); err != nil {
				return err
			}
			continue
		}

		// Create the file and write contents
		file, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, tr); err != nil {
			return err
		}
	}

	return nil
}
