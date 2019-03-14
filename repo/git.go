package repo

import (
	"bytes"
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	gitCTX = context.Background()
)

const (
	gitReposPerPage = 1000
)

func newAuthedGitClient(accesToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accesToken,
		},
	)
	tc := oauth2.NewClient(gitCTX, ts)
	return github.NewClient(tc)
}

func getAllGitReposForClient(client *github.Client) ([]*github.Repository, error) {
	// get all pages of results
	log.Infof("Getting all repos for client")
	var allRepos []*github.Repository
	repoListOptions := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: gitReposPerPage,
		},
	}
	for {
		repos, resp, err := client.Repositories.List(gitCTX, "", repoListOptions)
		if err != nil {
			// TODO not best implementation here -- will hit rate limiting, then early exit
			return allRepos, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		repoListOptions.Page = resp.NextPage
	}
	log.Debugf("Got %d repos", len(allRepos))
	return allRepos, nil
}

func getFileContents(client *github.Client, repo *github.Repository, path string) (*string, error) {
	ownerName := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	reader, err := client.Repositories.DownloadContents(gitCTX, ownerName, repoName, path, nil)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	contents := buf.String()
	return &contents, nil
}

func listDirectoryContents(client *github.Client, repo *github.Repository, root string) map[string]struct{} {
	ownerName := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	_, contents, _, _ := client.Repositories.GetContents(gitCTX, ownerName, repoName, root, nil)
	//FIXME no err handling

	files := make(map[string]struct{})
	for _, content := range contents {
		if content.GetType() == "file" {
			files[content.GetName()] = struct{}{}
		}
	}
	return files
}
