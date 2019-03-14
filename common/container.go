package common

import (
	"fmt"

	"github.com/google/go-github/github"
)

//Container links a Github repo with a generated Dockerfile
type Container struct {
	repo       *github.Repository
	dockerFile *string
}

func (c Container) String() string {
	var hasDockerStringPrefix string
	hasDocker := c.dockerFile != nil
	if !hasDocker {
		hasDockerStringPrefix = "No"
	} else {
		hasDockerStringPrefix = "Yes"
	}

	return fmt.Sprintf("{ [%s/%s], Docker: %s }", c.repo.GetOwner().GetLogin(), c.repo.GetName(), hasDockerStringPrefix)
}

//HasDockerFile returns whether or not the dockerfile property of the Container is nil
func (c Container) HasDockerFile() bool {
	return c.dockerFile != nil
}

//NewContainer constructs a new Container with the specified repo and dockerFile
func NewContainer(repo *github.Repository, dockerFile *string) Container {
	return Container{repo, dockerFile}
}
