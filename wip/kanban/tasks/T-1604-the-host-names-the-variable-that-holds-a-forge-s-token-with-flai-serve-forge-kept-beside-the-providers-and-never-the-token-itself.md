---
id: T-1604
type: task
nature: feature
title: The host names the variable that holds a forge's token with flai serve forge, kept beside the providers and never the token itself
status: backlog
parent: S-0392
owner: alex
created: 2026-10-16T04:33:32Z
updated: 2026-10-16T04:33:32Z
transitions: []
stream: S-0392
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
---
# T-1604 The host names the variable that holds a forge's token with flai serve forge, kept beside the providers and never the token itself

## Work

- In `flai/internal/config/config.go`, add `Forges`, a map from forge name to an entry holding `token_env` and an optional `api_url`, beside `Agent.Providers` (ADR-0129). It never holds a token.
- When no entry names the variable, `github` reads `GITHUB_TOKEN`.
- In `flai/cmd/serve_actions.go`, add `flai serve forge <name> [--token-env NAME] [--api-url URL] [--clear]`, built as `flai serve agent provider` is. It prints the entry as text, or as JSON with `--json`, naming the variable and whether it is set in this environment, never its value.
- Add a helper that resolves a forge's token from the environment at the moment of a call, so that callers never keep it.
- Test the config's round trip, the default name, and the command's output with the variable set and unset. Check that the value never appears in that output.
- First in its layer: it waits for no task.

## Done when

`flai serve forge github --token-env MY_TOKEN` records the name, `flai serve forge github` says whether the variable is set without printing it, and the tests pass.

## Notes
