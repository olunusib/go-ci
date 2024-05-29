package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/olunusib/go-ci/internal/config"
)

func StartServer(cfg *config.Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, cfg)
	})

	mux.HandleFunc("/logs/", logsHandler)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	maxRequestsPerSecond := 10
	burst := 20

	rateLimitedMux := RateLimit(maxRequestsPerSecond, burst, mux)

	log.Fatal(http.ListenAndServe(":"+cfg.PORT, rateLimitedMux))
}

type WebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		FullName string `json:"full_name"`
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
	HeadCommit struct {
		ID string `json:"id"`
	} `json:"head_commit"`
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

	go processWebhook(payload, cfg)

	response := map[string]string{"message": "Webhook received"}
	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	logID := r.URL.Path[len("/logs/"):]
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Error getting user home directory", http.StatusInternalServerError)
		return
	}
	logFilePath := filepath.Join(homeDir, "ci-logs", logID+".log")
	logFile, err := os.ReadFile(logFilePath)
	if err != nil {
		http.Error(w, "Log file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(logFile)
}
