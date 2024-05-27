package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	gitHTTP "github.com/go-git/go-git/v5/plumbing/transport/http"

	git "github.com/go-git/go-git/v5"
	"github.com/olunusib/go-ci/internal/ci"
	"github.com/olunusib/go-ci/internal/config"
)

type WebhookPayload struct {
	Repository struct {
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
}

func StartServer(cfg *config.Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, cfg)
	})

	maxRequestsPerSecond := 10
	burst := 20

	rateLimitedMux := RateLimit(maxRequestsPerSecond, burst, mux)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, rateLimitedMux))
}

func webhookHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		log.Printf("Error parsing request body: %v", err)
		return
	}

	go func() {
		processWebhook(payload, cfg)
	}()

	w.WriteHeader(http.StatusOK)
}

func cloneRepository(url, path, token string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &gitHTTP.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
		URL:      url,
		Depth:    1,
		Progress: os.Stdout,
	})
	return err
}

func processWebhook(payload WebhookPayload, cfg *config.Config) {
	repoURL := payload.Repository.CloneURL
	repoPath, err := os.MkdirTemp("/tmp", "")
	if err != nil {
		log.Printf("Error while creating temp directory: %v", err)
		return
	}

	pipelineFilePath := filepath.Join(repoPath, "ci", "pipeline.yaml")

	log.Printf("Cloning the repo: %s", repoURL)

	err = cloneRepository(repoURL, repoPath, cfg.GITHUB_TOKEN)
	if err != nil {
		log.Printf("Error while cloning repo: %v", err)
		return
	}

	err = os.Chdir(repoPath)
	if err != nil {
		log.Printf("Error while switching working dir: %v", err)
		return
	}

	pipelineConfig, err := ci.LoadConfig(pipelineFilePath)
	if err != nil {
		log.Printf("Error while loading pipeline config: %v", err)
		return
	}

	err = ci.ExecutePipeline(pipelineConfig)
	if err != nil {
		log.Printf("Pipeline execution failed: %v", err)
		return
	}
}
