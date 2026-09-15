.PHONY: check lint-md board stats dashboard dashboard-stop template-test flai flai-test flai-snapshot

check: ## Validate this repo against the standard
	flai check --strict

lint-md: ## Lint markdown
	npx --yes markdownlint-cli2 "**/*.md" "!**/node_modules/**"

board:
	flai board

stats:
	flai stats

dashboard:
	flai dashboard

dashboard-stop:
	flai dashboard stop

flai: ## Build flai into ./bin
	cd flai && go build -o ../bin/flai .

flai-test: ## Lint and race-test flai
	cd flai && golangci-lint run ./... && go test -race ./...

flai-snapshot: ## GoReleaser snapshot build of flai into flai/dist
	cd flai && goreleaser release --snapshot --clean --skip=publish

template-test: ## Render ./template into a temp dir and check it
	rm -rf /tmp/system-flow-template-test && flai new /tmp/system-flow-template-test --template ./template --defaults && flai check /tmp/system-flow-template-test --strict

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
