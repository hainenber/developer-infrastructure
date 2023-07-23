package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

type JcascJobConfig struct {
	Name   string
	Script string
}

var (
	jcascJobConfigTemplateString = `pipelineJob('{{.Name}}') {
	definition {
		cps {
			script("""
			{{.Script}}
		""".stripIndent())
		}
	}
}`
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Get Jenkinsfiles in ../jenkins/jobs directory
	jenkinsJobRootPath := filepath.Join(wd, "../jenkins/jobs")
	jenkinsJobs, err := os.ReadDir(jenkinsJobRootPath)
	if err != nil {
		panic(err)
	}

	jcascJobConfigTemplate := template.New("jcasc-config")
	jcascJobConfigTemplate = template.Must(jcascJobConfigTemplate.Parse(jcascJobConfigTemplateString))

	// Build Jenkinfile's content as JCasC's "jobs" pipeline
	var (
		jenkinsJobList []map[string]interface{}
	)
	for _, jenkinsJob := range jenkinsJobs {
		jenkinsJobPath := jenkinsJob.Name()
		jenkinsJobFiles, err := os.ReadDir(filepath.Join(jenkinsJobRootPath, jenkinsJobPath))
		if err != nil {
			panic(err)
		}
		for _, file := range jenkinsJobFiles {
			if file.Name() == "Jenkinsfile" {
				jenkinsFileContent, err := os.ReadFile(filepath.Join(jenkinsJobRootPath, jenkinsJobPath, "Jenkinsfile"))
				if err != nil {
					panic(err)
				}
				var buf bytes.Buffer
				if err := jcascJobConfigTemplate.Execute(&buf, JcascJobConfig{
					Name:   jenkinsJobPath,
					Script: string(jenkinsFileContent),
				}); err != nil {
					panic(err)
				}
				jenkinsJobList = append(jenkinsJobList, map[string]interface{}{
					"script": buf.String(),
				})
			}
		}
	}

	// Read the ../jenkins.yaml.scaffold as YAML format
	scaffoldJcascConfig, err := os.ReadFile(filepath.Join(wd, "../jenkins/casc-configs/jenkins.yaml.scaffold"))
	if err != nil {
		panic(err)
	}

	generatedJcascConfig := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(scaffoldJcascConfig, &generatedJcascConfig); err != nil {
		panic(err)
	}

	generatedJcascConfig["jobs"] = jenkinsJobList

	// Generate ../jenkins/casc-configs/jenkins.yaml from skeleton file ../jenkins/casc-configs/jenkins.yaml.scaffold
	marshalledGeneratedJcascConfig, err := yaml.Marshal(generatedJcascConfig)
	if err != nil {
		panic(err)
	}
	os.WriteFile(filepath.Join(wd, "../jenkins/casc-configs/jenkins.yaml"), marshalledGeneratedJcascConfig, 0644)
}
