.PHONY: check test integration smoke board stats dashboard dashboard-stop lint-md

check: ## Validate the repo against the system-flow standard
	scripts/check.sh

test: ## Behavior tests, fast, run on every iteration
	scripts/test.sh

integration: ## Integration tests, after test passes
	scripts/integration.sh

smoke: ## Smoke tests, after integration passes
	scripts/smoke.sh

lint-md: ## Lint markdown
	npx --yes markdownlint-cli2 "**/*.md" "!**/node_modules/**"

board: ## Print the kanban board
	flai board

stats: ## Print flow metrics
	flai stats

dashboard: ## Run the flaiover dashboard against this repo
	flai dashboard

dashboard-stop:
	flai dashboard stop

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
