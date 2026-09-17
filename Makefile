.PHONY: check flai test integration smoke flai-test flai-snapshot template-test install-tools lint-md board stats dashboard dashboard-stop flaiover-install flaiover-dev flaiover-build flaiover-test flaiover-image help

check: ## Validate this repo against the standard (flai check --strict)
	scripts/check.sh

flai: ## Build bin/flai from source
	scripts/flai-build.sh

test: ## Behavior tests (fast, run on every iteration)
	scripts/test.sh

integration: ## Integration tests (real git, monorepo round-trip), after test
	scripts/integration.sh

smoke: ## Smoke tests (template render and check, repo check, markdown lint), after integration
	scripts/smoke.sh

lint-md: ## Lint all markdown with the CI globs and config
	scripts/lint-md.sh

flai-test: ## Lint plus all three test tiers in order
	scripts/flai-test.sh

flai-snapshot: ## GoReleaser snapshot build into flai/dist
	scripts/flai-snapshot.sh

template-test: ## Render ./template and check the result
	scripts/template-test.sh

install-tools: ## Install golangci-lint v2, GoReleaser, and pnpm into the repo
	scripts/install-tools.sh

flaiover-install: ## pnpm install for flaiover
	scripts/flaiover-install.sh

flaiover-dev: ## flaiover dev server against this repo
	scripts/flaiover-dev.sh

flaiover-build: ## flaiover production build
	scripts/flaiover-build.sh

flaiover-test: ## flaiover lint, type check, unit tests
	scripts/flaiover-test.sh

flaiover-image: ## Build the flaiover image as flaiover:local
	scripts/flaiover-image.sh

lint-md: ## Lint markdown
	npx --yes markdownlint-cli2 "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**"

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
