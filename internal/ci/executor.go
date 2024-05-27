package ci

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func ExecutePipeline(pipelineConfig *Config) error {
	log.Printf("Starting to work on: %s", pipelineConfig.Name)

	for key, value := range pipelineConfig.Env {
		os.Setenv(key, value)
	}

	for _, step := range pipelineConfig.Steps {
		if err := runStep(step, pipelineConfig.Env); err != nil {
			return err
		}
	}

	log.Printf("CI pipeline completed successfully")
	return nil
}

func runStep(step Step, globalEnv map[string]string) error {
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
