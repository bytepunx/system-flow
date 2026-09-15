.PHONY: check lint-md board stats dashboard dashboard-stop template-test

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

template-test: ## Render ./template into a temp dir and check it
	rm -rf /tmp/system-flow-template-test && flai new /tmp/system-flow-template-test --template ./template --defaults && flai check /tmp/system-flow-template-test --strict

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
