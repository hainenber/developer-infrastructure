package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DEPLOY_MODE = "k8s"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Get Jenkinsfiles in ../jenkins/jobs directory
	// Populate their content as JCasC's "jobs" pipeline into template ../../jenkins/casc-configs/jcasc.yaml.scaffold
	var (
		jenkinsJobList       []map[string]interface{}
		jenkinsJobRootPath   = filepath.Join(wd, "../../jenkins/jobs")
		jcascDirPath         = filepath.Join(wd, "../../jenkins/casc-configs")
		generatedJcascConfig = make(map[interface{}]interface{})
	)
	// Visit child directory and extract Jenkinsfile's content for each defined job
	if err := filepath.WalkDir(jenkinsJobRootPath, func(path string, file fs.DirEntry, err error) error {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "Jenkinsfile") {
			jenkinsFileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			jenkinsJobList = append(jenkinsJobList, map[string]interface{}{
				"script": string(jenkinsFileContent),
			})
		}
		return nil
	}); err != nil {
		panic(err)
	}

	// Read the ../jcasc.yaml.scaffold in YAML format
	scaffoldJcascConfig, err := os.ReadFile(filepath.Join(jcascDirPath, fmt.Sprintf("jcasc.%s.yaml.scaffold", DEPLOY_MODE)))
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(scaffoldJcascConfig, &generatedJcascConfig); err != nil {
		panic(err)
	}

	// Generate ../jenkins/casc-configs/jcasc.yaml from skeleton file ../jenkins/casc-configs/jcasc.yaml.scaffold
	generatedJcascConfig["jobs"] = jenkinsJobList
	marshalledGeneratedJcascConfig, err := yaml.Marshal(generatedJcascConfig)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(jcascDirPath, "jcasc.yaml"), marshalledGeneratedJcascConfig, 0644); err != nil {
		panic(err)
	}
}
