.PHONY: check flai test integration smoke flai-test flai-snapshot template-test install-tools lint-md board stats dashboard dashboard-stop help

check: ## Validate this repo against the standard (flai check --strict)
	scripts/check.sh

flai: ## Build bin/flai from source
	scripts/flai-build.sh

test: ## Behavior tests (fast, run on every iteration)
	scripts/test.sh

integration: ## Integration tests (real git, monorepo round-trip), after test
	scripts/integration.sh

smoke: ## Smoke tests (template render and check, repo check), after integration
	scripts/smoke.sh

flai-test: ## Lint plus all three test tiers in order
	scripts/flai-test.sh

flai-snapshot: ## GoReleaser snapshot build into flai/dist
	scripts/flai-snapshot.sh

template-test: ## Render ./template and check the result
	scripts/template-test.sh

install-tools: ## Install golangci-lint v2 and GoReleaser into bin/
	scripts/install-tools.sh

lint-md: ## Lint markdown
	npx --yes markdownlint-cli2 "**/*.md" "!**/node_modules/**"

board: ## Print the kanban board
	scripts/flai.sh board

stats: ## Print flow metrics
	scripts/flai.sh stats

dashboard: ## Run the flaiover dashboard against this repo
	scripts/flai.sh dashboard

dashboard-stop:
	scripts/flai.sh dashboard stop

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
