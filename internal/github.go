package internal

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GitHubRelease struct {
	ZipballUrl string `json:"zipball_url"`
	Name       string `json:"name"`
}

func FetchLatestRelease(repo string, config *Config) {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	log.Printf("Fetching latest release from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Panic(err)
	}

	//get zipball
	resp, err = http.Get(release.ZipballUrl)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	//create file
	zipPath := config.GamePath + AddOnsFolder + "\\" + release.Name + ".zip"
	out, err := os.Create(zipPath)
	if err != nil {
		log.Panic(err)
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic(err)
	}

	unzippedPath := config.GamePath + AddOnsFolder + "\\" + release.Name
	if err = unzip(zipPath, unzippedPath); err != nil {
		log.Panic(err)
	}

}

func SplitProjectNameFromUrl(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

func unzip(zipPath, outDir string) error {
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
