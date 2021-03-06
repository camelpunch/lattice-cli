package docker_repository_name_formatter

import (
	"github.com/docker/docker/registry"
	"strings"
)

func FormatForReceptor(dockerRepositoryName string) (string, error) {

	_, err := registry.ParseRepositoryInfo(dockerRepositoryName)
	if err != nil {
		return "", err
	}

	if strings.Contains(dockerRepositoryName, "/") {
		return "docker:///" + dockerRepositoryName, nil
	} else {
		return "docker:///library/" + dockerRepositoryName, nil
	}
}
