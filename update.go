package main

import (
	"os"
	"io"
	"io/ioutil"
	"http"
	"fmt"
	"strings"
	"path"
	"exec"
)

/*
	Various URLs used for automatically updating.
*/
const (
	currentVersionUrl = "http://github.com/hgp/Go2Lunch/raw/master/VERSION"
	downloadUrl       = "http://github.com/downloads/hgp/Go2Lunch/lunch_"
	updateScriptUrl   = "http://github.com/hgp/Go2Lunch/raw/master/update_client.sh"
)

func getCurrentVersion() (string, os.Error) {
	res, _, err := http.Get(currentVersionUrl)
	if err != nil {
		return "", err
	}
	fmt.Println("1")
	if res.StatusCode != 200 {
		return "", os.NewError(fmt.Sprint("Could not check current version. Status code from server: ", res.StatusCode))
	}
	fmt.Println("2")
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("3")

	res.Body.Close()
	fmt.Println("4")
	return strings.TrimSpace(fmt.Sprintf("%s", buf)), nil
}

func needsUpdate() (result bool, currentVersion string, err os.Error) {
	currentVersion, err = getCurrentVersion()
	if err != nil {
		return
	}

	result = currentVersion != clientVersion
	return
}

func findLunch() (result string) {
	argPath := os.Args[0]
	result, _ = exec.LookPath(argPath)

	if !path.IsAbs(result) {
		// LookPath didn't find it, and the path isn't absolute, so put it relative to the current directory
		pwd, _ := os.Getwd()
		result = path.Join(pwd, path.Clean(argPath))
	}

	return
}

func downloadTempFile(url string, prefix string) (path string, err os.Error) {
	res, _, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", os.NewError("Could not download file. Status code from server: " + string(res.StatusCode))
	}

	dest, err := ioutil.TempFile("", prefix)
	if err != nil {
		return
	}
	defer dest.Close()

	_, err = io.Copy(dest, res.Body)
	if err != nil {
		return
	}

	return dest.Name(), nil
}

func downloadUpdateScript() (path string, err os.Error) {
	fmt.Println("downloadUpdateScript()")
	path, err = downloadTempFile(updateScriptUrl, "update_lunch")
	if err != nil {
		return
	}

	err = os.Chmod(path, 0700)
	return
}

func downloadNewVersion(version string) (path string, err os.Error) {
	return downloadTempFile(downloadUrl+version, "lunch")
}

func runUpdateScript(scriptPath, oldPath, newPath string) (err os.Error) {
	_, err = os.ForkExec(scriptPath, []string{scriptPath, oldPath, newPath}, os.Envs, "", []*os.File{os.Stdin, os.Stdout, os.Stderr})
	return
}

func CheckForUpdates() (errChan chan os.Error) {
	errChan = make(chan os.Error)
	go checkForUpdates(errChan)
	return errChan
}

func checkForUpdates(errChan chan os.Error) {
	update, version, err := needsUpdate()
	if err != nil {
		errChan <- err
	}

	if update {
		fmt.Println("An update is available. Would you like to download it? [Y/n] ")
		var result string
		fmt.Scanln(&result)

		fmt.Println()

		switch strings.ToLower(result) {
		case "", "y":
			fmt.Println("Downloading update script...")
			scriptPath, err := downloadUpdateScript()
			if err != nil {
				errChan <- err
			}
			fmt.Println("Update script downloaded to", scriptPath)
			fmt.Println()

			fmt.Println("Downloading new version...")
			newPath, err := downloadNewVersion(version)
			if err != nil {
				errChan <- err
			}
			fmt.Println("New version downloaded to", newPath)
			fmt.Println()

			fmt.Println("Running update script...")
			err = runUpdateScript(scriptPath, findLunch(), newPath)
			if err != nil {
				errChan <- err
			}

			// We should have forked the updater by now. We can exit.
			os.Exit(0)
		}
	}
	errChan <- nil
}
