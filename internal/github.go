package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GitHubRelease struct {
	ZipballUrl string `json:"zipball_url"`
	Name       string `json:"name"`
	AssetsUrl  string `json:"assets_url"`
}

type GithubReleaseAsset struct {
	DownloadUrl string `json:"browser_download_url"`
	Name        string `json:"name"`
}

func FetchLatestRelease(repo string, config *Config) {

	resp, err, release := getReleaseInfo(repo)
	resp, err, zipAsset := getZipReleaseAsset(resp, err, release)

	//download zip asset
	resp, err = http.Get(zipAsset.DownloadUrl)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	//create file
	zipPath, err := createZipFile(config, zipAsset, err, resp)
	defer CleanUpFile(zipPath)

	unzippedPath := config.GamePath + AddOnsFolder + "\\" + TrimFileExtension(zipAsset.Name, ".zip")
	err = Unzip(zipPath, unzippedPath)
	if err != nil {
		log.Panic(err)
	}
	defer CleanUpFile(unzippedPath)

	moveSubFilesAndUpdateConfig(repo, config, unzippedPath, release, err)
}

func getReleaseInfo(repo string) (*http.Response, error, GitHubRelease) {
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
	return resp, err, release
}

func getZipReleaseAsset(resp *http.Response, err error, release GitHubRelease) (*http.Response, error, GithubReleaseAsset) {
	resp, err = http.Get(release.AssetsUrl)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	var assets []GithubReleaseAsset
	if err = json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		log.Panic(err)
	}

	var zipAsset GithubReleaseAsset
	for _, asset := range assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			zipAsset = asset
		}
	}
	return resp, err, zipAsset
}

func SplitProjectNameFromUrl(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

func createZipFile(config *Config, zipAsset GithubReleaseAsset, err error, resp *http.Response) (string, error) {

	zipPath := config.GamePath + AddOnsFolder + "\\" + zipAsset.Name
	out, err := os.Create(zipPath)
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic(err)
	}

	return zipPath, err
}

func moveSubFilesAndUpdateConfig(repo string, config *Config, unzippedPath string, release GitHubRelease, err error) {
	folders := MoveFilesUpALevel(unzippedPath, config)
	repoUrl := "https://api.github.com/repos/" + repo
	version := release.Name
	config.AddAddOn(repoUrl, repo, version, folders)
	err = config.Save()
	if err != nil {
		log.Panic(err)
	}
}
