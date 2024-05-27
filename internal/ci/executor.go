package ci

import (
	"bytes"
	"log"
	"os/exec"
)

func ExecutePipeline(pipelineConfig *Config) error {
	log.Printf("Starting to work on: %s", pipelineConfig.Name)

	for _, step := range pipelineConfig.Steps {
		if err := runStep(step); err != nil {
			return err
		}
	}

	log.Printf("CI pipeline completed successfully")
	return nil
}

func runStep(step Step) error {
	log.Printf("Running step: %s", step.Name)
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
	return err
}
