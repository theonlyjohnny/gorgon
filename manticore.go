package main

import (
	"os"

	logPkg "github.com/theonlyjohnny/manticore/log"
	"github.com/theonlyjohnny/manticore/repo"
)

var (
	log       = logPkg.Log
	testToken = os.Getenv("GITHUB_TOKEN")
)

func main() {
	if testToken == "" {
		log.Errorf("No GITHUB_TOKEN env var set -- exiting")
		os.Exit(1)
	}

	containers, err := repo.GetContainersFromGitUser(testToken)
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		os.Exit(1)
	}
	log.Infof("Containers: %s", containers)
}
