package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type GoProxyHost struct {
	GoProxyHost string
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var (
		jcascDirPath           = filepath.Join(wd, "../../jenkins/casc-configs")
		athensContainerID      string
		athensPrivateIpAddress string
	)

	jcascTemplateContent, err := os.ReadFile(filepath.Join(jcascDirPath, "jenkins.yaml.scaffold"))
	if err != nil {
		panic(err)
	}

	jcascTemplate := template.Must(template.New("jcasc-config").Parse(string(jcascTemplateContent)))

	// Get IP of running container named "athens"
	colimaDockerSocketPath := filepath.Join(os.Getenv("HOME"), ".colima/default/docker.sock")
	if _, err := os.Stat(colimaDockerSocketPath); !os.IsExist(err) {
		if err := os.Setenv("DOCKER_HOST", fmt.Sprintf("unix://%s/.colima/default/docker.sock", os.Getenv("HOME"))); err != nil {
			panic(err)
		}
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	runningContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, ctn := range runningContainers {
		if len(ctn.Names) == 1 && ctn.Names[0] == "/athens" {
			athensContainerID = ctn.ID
		}
	}

	if athensContainerID == "" {
		panic(fmt.Errorf("cannot find Athens container ID. Exit"))
	}

	// Add Athens's IP  as Jenkins global variable
	athensInspectionResult, err := cli.ContainerInspect(ctx, athensContainerID)
	if err != nil {
		panic(err)
	}

	if len(athensInspectionResult.NetworkSettings.Networks) > 1 {
		panic(fmt.Errorf("athens container exists in >1 Docker network. Exit"))
	}
	for _, networkSetting := range athensInspectionResult.NetworkSettings.Networks {
		athensPrivateIpAddress = networkSetting.IPAddress
	}

	athensContainerData := GoProxyHost{
		GoProxyHost: fmt.Sprintf("http://%s:3000", athensPrivateIpAddress),
	}

	// Read the ../jenkins.yaml.scaffold in YAML format
	generatedJcascConfigFile, err := os.Create(filepath.Join(jcascDirPath, "jenkins.yaml"))
	if err != nil {
		panic(err)
	}
	defer generatedJcascConfigFile.Close()

	// Generate ../jenkins/casc-configs/jenkins.yaml from skeleton file ../jenkins/casc-configs/jenkins.yaml.scaffold
	jcascTemplate.Execute(generatedJcascConfigFile, athensContainerData)
}
