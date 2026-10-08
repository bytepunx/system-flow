---
id: T-1426
type: task
nature: improvement
title: The claude version check runs through the project's provider, records the provider per version, and skips a provider whose key is unset
status: backlog
parent: S-0360
owner: alex
created: 2026-10-08T08:59:47Z
updated: 2026-10-08T08:59:47Z
transitions: []
stream: S-0360
tags: [cli]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go]
---
# T-1426 The claude version check runs through the project's provider, records the provider per version, and skips a provider whose key is unset

## Work

- For each served project whose agents run `claude-code`, resolve the default agent's provider as S-0350's start resolves it. Build the check's command with the adapter's `Start` for that provider, so the environment and the `EnvFrom` copy are the same as a story's start.
- `claude-checks.json` keys each record by version and provider (`anthropic` when none); reading an older record without a provider counts it as `anthropic`.
- A provider whose key's variable is unset or empty is skipped with one line in `flai serve`'s log, and no thread is opened.
- Tests with a fake `claude` that prints the environment it saw: a check with a provider, one without, a version passed through Anthropic and checked again through a provider, and an unset key.

First layer.

## Done when

- `flai test flai/internal/serve/` passes.

## Notes

Layer 1 of S-0360.
