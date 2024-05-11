package internal

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func TrimFileExtension(filePath, extension string) string {
	return strings.TrimSuffix(filePath, extension)
}

func Unzip(zipPath, outDir string) error {
	// Open the zip file specified by zipPath.
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Ensure the output directory exists.
	os.MkdirAll(outDir, 0755)

	// Iterate through the files in the archive.
	for _, f := range r.File {
		// Create full path for the extracted file.
		filePath := filepath.Join(outDir, f.Name)

		// If the file is a directory, create it and move to the next file.
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// Create the directories leading up to the file if they don't exist.
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// Open the file within the zip archive.
		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create a new file within the output directory.
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Copy the file contents to the new file.
		if _, err := io.Copy(outFile, srcFile); err != nil {
			return err
		}
	}

	return nil
}

func MoveFilesUpALevel(filePath string, config *Config) (folders []string) {
	addOnsPath := config.GamePath + "\\" + AddOnsFolder

	items, err := os.ReadDir(filePath)
	if err != nil {
		log.Panic(err)
	}

	for _, item := range items {
		if item.IsDir() {
			dirPath := filepath.Join(filePath, item.Name())
			newPath := filepath.Join(addOnsPath, item.Name())

			err := os.Rename(dirPath, newPath)
			if err != nil {
				log.Panic(err)
			}
			folders = append(folders, newPath)
		}
	}

	return folders
}

func CleanUpFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func BackUpExistingFile(file string) string {
	backUpFileName := file + ".backup"

	err := os.Rename(file, backUpFileName)
	if err != nil {
		panic(err)
	}

	return backUpFileName
}

func RestoreBackUps(backedUpFiles []string) {
	for _, file := range backedUpFiles {
		err := os.Rename(file, strings.TrimSuffix(file, ".backup"))
		if err != nil {
			log.Println("Error restoring from backup: " + file)
			panic(err)
		}
	}
}
