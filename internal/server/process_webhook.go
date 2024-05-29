package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

	runID := generateRunID()
	logURL := fmt.Sprintf("%s/logs/%s", cfg.SERVER_BASE_URL, runID)

	fmt.Println(logURL)

	githubClient.SetCommitStatus(payload, "pending", "Processing your request", logURL)

	repoPath, err := createTempDir()
	if err != nil {
		log.Printf("Error while creating temp directory: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to create temp directory", logURL)
		return
	}
	defer os.RemoveAll(repoPath)

	branch := payload.Ref
	repoURL := payload.Repository.CloneURL
	log.Printf("Cloning the repo: %s on branch %s", repoURL, branch)

	if err := cloneRepository(repoURL, repoPath, cfg.GITHUB_TOKEN, branch); err != nil {
		log.Printf("Error while cloning repo: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to clone repository", logURL)
		return
	}

	if err := changeWorkingDir(repoPath); err != nil {
		log.Printf("Error while switching working dir: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to switch working directory", logURL)
		return
	}

	pipelineFilePath, err := getPipelineFilePath()
	if err != nil {
		log.Printf("No pipeline configuration file found")
		githubClient.SetCommitStatus(payload, "failure", "No pipeline configuration file found", logURL)
		return
	}

	if err := loadAndExecutePipeline(pipelineFilePath, runID); err != nil {
		log.Printf("Pipeline execution failed: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Pipeline execution failed", logURL)
		return
	}

	githubClient.SetCommitStatus(payload, "success", "Processing complete", logURL)
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

func loadAndExecutePipeline(pipelineFilePath, runID string) error {
	pipelineConfig, err := ci.LoadConfig(pipelineFilePath)
	if err != nil {
		return err
	}
	_, err = ci.ExecutePipeline(pipelineConfig, runID)
	return err
}

func generateRunID() string {
	return uuid.New().String()
}
