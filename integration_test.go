package main

import (
	"fmt"
	"github.com/epiphany-platform/m-azure-kubernetes-service/cmd"
	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
	"os"
	"testing"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		name               string
		wantOutputTemplate string
		parameters         map[string]string
	}{
		{
			name: "default metadata",
			wantOutputTemplate: `labels:
  kind: infrastructure
  name: Azure Kubernetes Service
  provider: azure
  provides-pubips: true
  short: azks
  version: %s
`,
			parameters: nil,
		},
		{
			name:               "json metadata",
			wantOutputTemplate: "{\"labels\":{\"kind\":\"infrastructure\",\"name\":\"Azure Kubernetes Service\",\"provider\":\"azure\",\"provides-pubips\":true,\"short\":\"azks\",\"version\":\"%s\"}}",
			parameters:         map[string]string{"--json": "true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := dockerRun(t, "metadata", tt.parameters, nil, "")
			if diff := deep.Equal(gotOutput, fmt.Sprintf(tt.wantOutputTemplate, cmd.Version)); diff != nil {
				t.Error(diff)
			}
		})
	}
}

// dockerRun function wraps docker run operation and returns `docker run` output.
func dockerRun(t *testing.T, command string, parameters map[string]string, environments map[string]string, sharedPath string) string {
	commandWithParameters := []string{command}
	for k, v := range parameters {
		commandWithParameters = append(commandWithParameters, fmt.Sprintf("%s=%s", k, v))
	}

	var opts *docker.RunOptions
	if sharedPath != "" {
		opts = &docker.RunOptions{
			Command: commandWithParameters,
			Remove:  true,
			Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
		}
	} else {
		opts = &docker.RunOptions{
			Command: commandWithParameters,
			Remove:  true,
		}
	}
	var envs []string
	for k, v := range environments {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	opts.EnvironmentVariables = envs

	//in case of error Run function calls FailNow anyways
	return docker.Run(t, fmt.Sprintf("%s:%s", prepareImageTag(t), cmd.Version), opts)
}

// prepareImageTag returns IMAGE_REPOSITORY environment variable
func prepareImageTag(t *testing.T) string {
	imageRepository := os.Getenv("IMAGE_REPOSITORY")
	if len(imageRepository) == 0 {
		t.Fatal("expected IMAGE_REPOSITORY environment variable")
	}
	return imageRepository
}
