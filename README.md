[![Deploy](https://github.com/olunusib/go-ci/actions/workflows/deploy.yaml/badge.svg)](https://github.com/olunusib/go-ci/actions/workflows/deploy.yaml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/olunusib/go-ci)
![GitHub issues](https://img.shields.io/github/issues/olunusib/go-ci)
![GitHub pull requests](https://img.shields.io/github/issues-pr/olunusib/go-ci)
![GitHub](https://img.shields.io/github/license/olunusib/go-ci)

# Go-CI

Go-CI is a self-hosted continuous integration (CI) server that listens for webhooks and executes CI pipelines in response to repository events.  

## Getting Started

### **1 Before running the server, set up the required environment variables:  

```bash
# Token with repo scope (for accessing repositories)
export GITHUB_TOKEN=<>
# Base URL where this CI server is hosted
export SERVER_BASE_URL=<>
```

### **2 Use Docker to pull the latest image and start the server:

```bash
docker pull ghcr.io/olunusib/go-ci:latest

docker run -p 8080:8080 \
  -e GITHUB_TOKEN=$GITHUB_TOKEN \
  -e SERVER_BASE_URL=$SERVER_BASE_URL \
  -d ghcr.io/olunusib/go-ci:latest
```

### **3 Set Up Webhooks:

To trigger builds, configure a GitHub Webhook in your repository with the following details:

```bash
Payload URL: http://your-server-url/webhook
Content type: application/json
Events: Push, Pull Requests, or any relevant event.
```

## Example GitHub CI pipeline

This should be places in the `/ci` directory, e.g., `/ci/pipeline.yaml`

```yaml
name: CI
env:
  GLOBAL_ENV1: value1
  GLOBAL_ENV2: value2
steps:
  - name: Install Python and pip
    command: |
      if ! command -v python3 &> /dev/null; then
        apt-get update && apt-get install -y python3 python3-venv
      fi
  - name: Set up virtual environment
    command: python3 -m venv venv
  - name: Install dependencies
    command: venv/bin/pip install -r requirements.txt
  - name: Run tests
    command: venv/bin/python -m unittest discover -s tests
  - name: Inspect Env Vars
    env:
      STEP_ENV: value3
    command: echo $GLOBAL_ENV1 $GLOBAL_ENV2 $STEP_ENV
```
