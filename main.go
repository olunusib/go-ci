package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/olunusib/go-ci/pipeline"
)

func main() {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("Environment variable GITHUB_TOKEN not set")
	}

	repoURL, ok := os.LookupEnv("SOURCE_URL")
	if !ok {
		log.Fatal("Environment variable SOURCE_URL not set")
	}

	repoPath, err := os.MkdirTemp("/tmp", "")
	if err != nil {
		log.Fatalf("Error while creating temp directory: %v", err)
	}

	pipelineFilePath := filepath.Join(repoPath, "ci", "pipeline.yaml")

	log.Printf("Cloning the repo: %s", repoURL)

	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
		URL:      repoURL,
		Depth:    1,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatalf("Error while cloning repo: %v", err)
	}

	err = os.Chdir(repoPath)
	if err != nil {
		log.Fatalf("Error while switching working dir: %v", err)
	}

	pipelineConfig, err := pipeline.LoadConfig(pipelineFilePath)
	if err != nil {
		log.Fatalf("Error while loading pipeline config: %v", err)
	}

	log.Printf("Starting to work on: %s", pipelineConfig.Name)

	for _, step := range pipelineConfig.Steps {
		log.Printf("Running step: %s", step.Name)
		cmd := exec.Command("sh", "-c", step.Command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			log.Fatalf("Step %s failed: %v", step.Name, err)
		}
	}

	log.Printf("CI pipeline completed successfully")
}
