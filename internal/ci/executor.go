package ci

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ExecutePipeline(pipelineConfig *Config, runID string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	logDir := filepath.Join(homeDir, "ci-logs")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.log", runID))
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return "", fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Printf("Starting to work on: %s", pipelineConfig.Name)

	for key, value := range pipelineConfig.Env {
		os.Setenv(key, value)
	}

	for _, step := range pipelineConfig.Steps {
		if err := runStep(step); err != nil {
			return runID, err
		}
	}

	log.Printf("CI pipeline completed successfully")
	return runID, nil
}

func runStep(step Step) error {
	log.Printf("Running step: %s", step.Name)

	for key, value := range step.Env {
		os.Setenv(key, value)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", step.Command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Printf("Output of step %s:\n%s", step.Name, stdout.String())
	if stderr.Len() > 0 {
		log.Printf("Error output of step %s:\n%s", step.Name, stderr.String())
	}
	if err != nil {
		log.Printf("Step %s failed: %v", step.Name, err)
	}

	for key := range step.Env {
		os.Unsetenv(key)
	}

	return err
}
