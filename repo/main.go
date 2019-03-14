package repo

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/theonlyjohnny/manticore/common"
	logPkg "github.com/theonlyjohnny/manticore/log"
)

var (
	log        = logPkg.Log
	generators = map[string]dockerfileGenerator{
		"package.json": generateNodeDockerfile,
	}
)

type dockerfileGenerator func(string) *string

//GetContainersFromGitUser returns all the Git Repos associated to the provided accesToken
//in the form of Containers
func GetContainersFromGitUser(accessToken string) ([]common.Container, error) {
	client := newAuthedGitClient(accessToken)
	return getContainers(client)
}

func getContainers(client *github.Client) ([]common.Container, error) {
	var allGitRepos []*github.Repository
	var containers []common.Container
	var operated int
	var err error

	allGitRepos, err = getAllGitReposForClient(client)
	if err != nil {
		return containers, err
	}

	repoCount := len(allGitRepos)

	containersChan := make(chan common.Container, repoCount)

	for _, repo := range allGitRepos {
		go getContainerFromGitRepo(containersChan, client, repo)
	}

	for container := range containersChan {
		if container.HasDockerFile() {
			containers = append(containers, container)
		}
		operated++
		log.Debugf("Operated on %d/%d", operated, repoCount)
		if operated == repoCount {
			close(containersChan)
		}
	}

	return containers, nil
}

func getContainerFromGitRepo(output chan<- common.Container, client *github.Client, repo *github.Repository) {

	var dockerFile *string

	dockerFile = getDockerFileFromRepo(client, repo)

	output <- common.NewContainer(repo, dockerFile)
}

func checkDockerfileGenerator(client *github.Client, repo *github.Repository, files map[string]struct{}, generatorName string, generatorFunc dockerfileGenerator) *string {
	ownerName := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	prefix := fmt.Sprintf("%s/%s:", ownerName, repoName)
	if _, ok := files[generatorName]; ok {
		file, err := getFileContents(client, repo, generatorName)
		if err != nil {
			log.Errorf("%s Has %s in toplevel, but was unable to download contents: %s", prefix, generatorName, err.Error())
			return nil
		}

		return generatorFunc(*file)
	}

	return nil
}

func getDockerFileFromRepo(client *github.Client, repo *github.Repository) *string {
	var dockerfile *string

	files := listDirectoryContents(client, repo, "")
	for generatorName, generatorFunc := range generators {
		generated := checkDockerfileGenerator(client, repo, files, generatorName, generatorFunc)
		if generated != nil {
			dockerfile = generated
			break
		}
	}

	return dockerfile
}
