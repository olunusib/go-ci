# go-ci

```bash
# Token with repo scope
export GITHUB_TOKEN=<>
# Where this is being hosted
export SERVER_BASE_URL=<>

# Pull image and start server
docker pull ghcr.io/olunusib/go-ci:latest

docker run -p 8080:8080 -e GITHUB_TOKEN=$GITHUB_TOKEN -e SERVER_BASE_URL=$SERVER_BASE_URL -d ghcr.io/olunusib/go-ci:latest
```
