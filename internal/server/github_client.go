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

	status, err := client.constructStatusPayload(state, description, targetURL)
	if err != nil {
		log.Printf("Error constructing status payload: %v", err)
		return
	}

	url := client.constructStatusURL(payload)
	req, err := client.constructStatusRequest(url, status)
	if err != nil {
		log.Printf("Error constructing status request: %v", err)
		return
	}

	client.sendStatusRequest(req)
}

func (client *GitHubClient) constructStatusPayload(state, description, targetURL string) ([]byte, error) {
	status := GitHubStatus{
		State:       state,
		Description: description,
		TargetURL:   targetURL,
		Context:     "continuous-integration/go-ci",
	}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	return statusJSON, nil
}

func (client *GitHubClient) constructStatusURL(payload WebhookPayload) string {
	repo := payload.Repository.FullName
	commitSHA := payload.HeadCommit.ID
	return fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s", repo, commitSHA)
}

func (client *GitHubClient) constructStatusRequest(url string, statusJSON []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(statusJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+client.Token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (client *GitHubClient) sendStatusRequest(req *http.Request) {
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
