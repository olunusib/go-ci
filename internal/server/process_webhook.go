package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/olunusib/go-ci/internal/ci"
	"github.com/olunusib/go-ci/internal/config"

	gitHTTP "github.com/go-git/go-git/v5/plumbing/transport/http"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func processWebhook(payload WebhookPayload, cfg *config.Config) {
	githubClient := &GitHubClient{
		Token: cfg.GITHUB_TOKEN,
	}

	// this should be the actual URL of your CI/CD job page
	targetURL := "https://go-ci.example.com"

	githubClient.SetCommitStatus(payload, "pending", "Processing your request", targetURL)

	repoURL := payload.Repository.CloneURL
	repoPath, err := os.MkdirTemp("/tmp", "")
	if err != nil {
		log.Printf("Error while creating temp directory: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to create temp directory", targetURL)
		return
	}

	defer os.RemoveAll(repoPath)

	pipelineFilePath := filepath.Join(repoPath, "ci", "pipeline.yaml")

	branch := payload.Ref

	log.Printf("Cloning the repo: %s on branch %s", repoURL, branch)

	err = cloneRepository(repoURL, repoPath, cfg.GITHUB_TOKEN, branch)
	if err != nil {
		log.Printf("Error while cloning repo: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to clone repository", targetURL)
		return
	}

	err = os.Chdir(repoPath)
	if err != nil {
		log.Printf("Error while switching working dir: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to switch working directory", targetURL)
		return
	}

	pipelineConfig, err := ci.LoadConfig(pipelineFilePath)
	if err != nil {
		log.Printf("Error while loading pipeline config: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Failed to load pipeline config", targetURL)
		return
	}

	err = ci.ExecutePipeline(pipelineConfig)
	if err != nil {
		log.Printf("Pipeline execution failed: %v", err)
		githubClient.SetCommitStatus(payload, "failure", "Pipeline execution failed", targetURL)
		return
	}

	githubClient.SetCommitStatus(payload, "success", "Processing complete", targetURL)
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
