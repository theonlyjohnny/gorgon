# gorgon [WIP]

### Description
A Golang tool that scrapes your Github profile (through [Github Access Tokens](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)) and converts Github repos into Docker containers.
Gorgon "freezes" your Github repos into a Docker container from whenever its run. It uses your own Dockerfile if present, or tries to generate one for you based on inferred project architecture.

The idea is that this could be used as a step in a custom CI pipeline, with every commit creating a new docker repo automatically

### TODO:
 - [ ] Auth w/ DockerHub
 - [ ] Check if container is already pushed up
 - [ ] Webhook for post-deployment?
 - [ ] Add more Dockerfile templates than NodeJS
 - [ ] Docker build security
