package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/olunusib/go-ci/internal/ci"
	"github.com/olunusib/go-ci/internal/config"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitHTTP "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func processWebhook(payload WebhookPayload, cfg *config.Config) {
	githubClient := &GitHubClient{
		Token: cfg.GITHUB_TOKEN,
	}

	// This should be the actual URL of your CI/CD job page
	targetURL := "https://go-ci.example.com"

	githubClient.SetCommitStatus(payload, "pending", "Processing your request", targetURL)

	repoPath, err := createTempDir()
	if err != nil {
		log.Printf("Error while creating temp directory: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to create temp directory", targetURL)
		return
	}
	defer os.RemoveAll(repoPath)

	branch := payload.Ref
	repoURL := payload.Repository.CloneURL
	log.Printf("Cloning the repo: %s on branch %s", repoURL, branch)

	if err := cloneRepository(repoURL, repoPath, cfg.GITHUB_TOKEN, branch); err != nil {
		log.Printf("Error while cloning repo: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to clone repository", targetURL)
		return
	}

	if err := changeWorkingDir(repoPath); err != nil {
		log.Printf("Error while switching working dir: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to switch working directory", targetURL)
		return
	}

	pipelineFilePath, err := getPipelineFilePath()
	if err != nil {
		log.Printf("No pipeline configuration file found")
		return
	}

	if err := loadAndExecutePipeline(pipelineFilePath); err != nil {
		log.Printf("Pipeline execution failed: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Pipeline execution failed", targetURL)
		return
	}

	githubClient.SetCommitStatus(payload, "success", "Processing complete", targetURL)
}

func createTempDir() (string, error) {
	return os.MkdirTemp("/tmp", "")
}

func getPipelineFilePath() (string, error) {
	possiblePaths := []string{
		filepath.Join("ci", "pipeline.yml"),
		filepath.Join("ci", "pipeline.yaml"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}

func cloneRepository(url, path, token, branch string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &gitHTTP.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
		URL:           url,
		ReferenceName: plumbing.ReferenceName(branch),
		Depth:         1,
		Progress:      os.Stdout,
	})
	return err
}

func changeWorkingDir(path string) error {
	return os.Chdir(path)
}

func loadAndExecutePipeline(pipelineFilePath string) error {
	pipelineConfig, err := ci.LoadConfig(pipelineFilePath)
	if err != nil {
		return err
	}
	return ci.ExecutePipeline(pipelineConfig)
}
