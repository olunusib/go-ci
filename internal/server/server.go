package server

import (
	"log"
	"net/http"

	"github.com/olunusib/go-ci/internal/config"
)

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
