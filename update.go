package main

import (
	"os"
	"io/ioutil"
	"http"
	"fmt"
	"strings"
)

/*
	Various URLs used for automatically updating.
*/
const (
	currentVersionUrl = "http://github.com/hgp/Go2Lunch/raw/master/VERSION"
	downloadUrl = "http://github.com/downloads/hgp/Go2Lunch/lunch_"
)

func getCurrentVersion() (string, os.Error) {
	res, _, err := http.Get(currentVersionUrl)
	if err != nil {
		return "", err
	}
	
	if res.StatusCode != 200 {
		return "", os.NewError(fmt.Sprint("Could not check current version. Status code from server: ", res.StatusCode))
	}
	
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	
	res.Body.Close()
	return strings.TrimSpace(fmt.Sprintf("%s", buf)), nil
}

func needsUpdate() (bool, os.Error) {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		return false, err
	}
	
	return currentVersion != clientVersion, nil
}

func CheckForUpdates() (err os.Error) {
	update, err := needsUpdate()
	if err != nil {
		return
	}
	
	if update {
		fmt.Print("An update is available. Would you like to download it? [Y/n] ")
		var result string
		fmt.Scanln(&result)
		
		fmt.Println()
		
		switch strings.ToLower(result) {
		case "", "y":
			version, _ := getCurrentVersion()
			fmt.Println("Automatic downloading of updates is not available yet.")
			fmt.Println("You can download the new version at", downloadUrl + version)
			fmt.Println()
			return
		}
		
	}
	return
}
