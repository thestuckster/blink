package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitHubRelease struct {
	ZipballUrl string `json:"zipball_url"`
	Name       string `json:"name"`
	AssetsUrl  string `json:"assets_url"`
}

type GithubReleaseAsset struct {
	DownloadUrl string `json:"browser_download_url"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func DownloadLatestRelease(repo string, config *Config) {

	resp, err, release := getReleaseInfo(repo)
	if err != nil {
		log.Panic(err)
	}
	resp, err, zipAsset := getZipReleaseAsset(resp, err, *release)
	if err != nil {
		log.Panic(err)
	}

	//download zip asset
	resp, err = http.Get(zipAsset.DownloadUrl)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	//create file
	zipPath, err := createZipFile(config, *zipAsset, resp)

	unzippedPath := config.GamePath + AddOnsFolder + "\\" + TrimFileExtension(zipAsset.Name, ".zip")
	err = Unzip(zipPath, unzippedPath)
	if err != nil {
		log.Panic(err)
	}

	moveSubFilesAndSaveAddOnDetails(repo, config, unzippedPath, *release, *zipAsset)

	err = CleanUpFile(zipPath)
	if err != nil {
		log.Panic(err)
	}

	err = CleanUpFile(unzippedPath)
	if err != nil {
		log.Panic(err)
	}

}

func SplitProjectNameFromUrl(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

func Update(url string, addOnDetails *AddOn, config *Config) (bool, error) {
	resp, err, releaseInfo := getReleaseInfo(addOnDetails.Repo)
	if err != nil {
		return false, err
	}

	resp, err, zipAsset := getZipReleaseAsset(resp, err, *releaseInfo)
	if err != nil {
		return false, err
	}

	assetCreatedDate, err := time.Parse(time.RFC3339, zipAsset.CreatedAt)
	if err != nil {
		return false, err
	}

	if !addOnDetails.CreatedAt.Before(assetCreatedDate) {
		log.Println(addOnDetails.Repo + " is already update to date")
		return false, nil
	}

	zipPath, err := createZipFile(config, *zipAsset, resp)
	if err != nil {
		return false, err
	}

	unzippedPath := config.GamePath + AddOnsFolder + "\\" + TrimFileExtension(zipAsset.Name, ".zip")
	err = Unzip(zipPath, unzippedPath)
	if err != nil {
		return false, err
	}

	err = moveSubFilesAndUpdateConfig(addOnDetails.Repo, config, unzippedPath, releaseInfo, zipAsset)
	if err != nil {
		return false, err
	}

	err = CleanUpFile(zipPath)
	if err != nil {
		return false, err
	}

	err = CleanUpFile(unzippedPath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getReleaseInfo(repo string) (*http.Response, error, *GitHubRelease) {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	log.Printf("Fetching latest release from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err, nil
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err, nil
	}
	return resp, err, &release
}

func getZipReleaseAsset(resp *http.Response, err error, release GitHubRelease) (*http.Response, error, *GithubReleaseAsset) {
	resp, err = http.Get(release.AssetsUrl)
	if err != nil {
		return nil, err, nil
	}
	defer resp.Body.Close()

	var assets []GithubReleaseAsset
	if err = json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return nil, err, nil
	}

	var zipAsset GithubReleaseAsset
	for _, asset := range assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			zipAsset = asset
		}
	}
	return resp, err, &zipAsset
}

func createZipFile(config *Config, zipAsset GithubReleaseAsset, resp *http.Response) (string, error) {

	zipPath := config.GamePath + AddOnsFolder + "\\" + zipAsset.Name
	out, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}

	//TODO: somehow this file is already closed???
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	out.Close()
	return zipPath, err
}

func moveSubFilesAndSaveAddOnDetails(repo string,
	config *Config,
	unzippedPath string,
	release GitHubRelease,
	asset GithubReleaseAsset) {

	folders := MoveFilesUpALevel(unzippedPath, config)
	repoUrl := "https://api.github.com/repos/" + repo
	version := release.Name
	createdAt, err := time.Parse(time.RFC3339, asset.CreatedAt)
	if err != nil {
		log.Panic(err)
	}

	config.AddAddOn(repoUrl, repo, version, createdAt, folders)
	err = config.Save()
	if err != nil {
		log.Panic(err)
	}
}

func moveSubFilesAndUpdateConfig(
	repo string,
	config *Config,
	unzippedPath string,
	release *GitHubRelease,
	asset *GithubReleaseAsset) error {

	folders := MoveFilesUpALevel(unzippedPath, config)
	repoUrl := "https://api.github.com/repos/" + repo
	version := release.Name
	createdAt, err := time.Parse(time.RFC3339, asset.CreatedAt)
	if err != nil {
		return err
	}

	i, _ := FindAddOnDetails(repo, config)
	newAddOnDetails := AddOn{
		Url:       repoUrl,
		Repo:      repo,
		Version:   version,
		Folders:   folders,
		CreatedAt: createdAt,
	}

	config.AddOns[i] = newAddOnDetails
	err = config.Save()
	if err != nil {
		return err
	}

	return nil
}
