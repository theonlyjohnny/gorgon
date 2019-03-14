package main

import (
	"os"

	"github.com/schollz/progressbar"
	"github.com/theonlyjohnny/manticore/common"
	logPkg "github.com/theonlyjohnny/manticore/log"
	"github.com/theonlyjohnny/manticore/repo"
)

var (
	bar       *progressbar.ProgressBar
	log       = logPkg.Log
	testToken = os.Getenv("GITHUB_TOKEN")
)

func listenToProgressChan(progressChan <-chan common.Progress) {
	for {
		progress := <-progressChan
		if bar == nil {
			bar = progressbar.New(progress.Total)
			bar.RenderBlank()
		}
		bar.Set(progress.Current)
	}
}

func main() {
	if testToken == "" {
		log.Errorf("No GITHUB_TOKEN env var set -- exiting")
		os.Exit(1)
	}

	progressChan := make(chan common.Progress, 1)

	go listenToProgressChan(progressChan)

	containers, err := repo.GetContainersFromGitUser(testToken, progressChan)
	close(progressChan)

	if err != nil {
		log.Errorf("Error: %s", err.Error())
		os.Exit(1)
	}
	log.Infof("Containers: %s", containers)
}
