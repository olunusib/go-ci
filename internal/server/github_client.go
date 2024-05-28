package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GitHubStatus struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
	Context     string `json:"context,omitempty"`
}

type GitHubClient struct {
	Token string
}

func (client *GitHubClient) SetCommitStatus(payload WebhookPayload, state, description, targetURL string) {
	if client.Token == "" {
		log.Println("GitHub token is not set")
		return
	}

	status := GitHubStatus{
		State:       state,
		Description: description,
		TargetURL:   targetURL,
		Context:     "continuous-integration/go-ci",
	}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		log.Printf("Error marshaling status JSON: %v", err)
		return
	}

	repo := payload.Repository.FullName
	commitSHA := payload.HeadCommit.ID

	log.Printf("Setting commit status for %s", commitSHA)

	url := fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s", repo, commitSHA)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(statusJSON))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Authorization", "token "+client.Token)
	req.Header.Set("Content-Type", "application/json")

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to set status: %s", resp.Status)
	}
}
