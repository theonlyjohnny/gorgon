package repo

import (
	"encoding/json"
	"fmt"
)

type packageJSON struct {
	Entrypoint *string           `json:"main"`
	Scripts    map[string]string `json:"scripts"`
}

func generateNodeDockerfile(packageJSONContents string) *string {
	var packageJSON packageJSON
	var runCmd string

	err := json.Unmarshal([]byte(packageJSONContents), &packageJSON)
	if err != nil {
		return nil
	}
	if startScript, ok := packageJSON.Scripts["start"]; ok {
		runCmd = startScript
	} else if packageJSON.Entrypoint != nil {
		runCmd = *packageJSON.Entrypoint
	}

	if runCmd == "" {
		return nil
	}

	contents := fmt.Sprintf(
		`FROM node:11.11.0-alpine
	COPY package.json package.json
	RUN yarn i
	COPY . .
	ENTRYPOINT node
	CMD %s`, runCmd)

	return &contents
}
