---
id: T-1707
type: task
nature: feature
title: The opencode configuration carries a provider entry derived from the host's providers map, with the key named and never written, and refuses an api OpenCode cannot speak
status: backlog
parent: S-0410
owner: alex
created: 2026-10-16T05:18:48Z
updated: 2026-10-16T05:18:48Z
transitions: []
stream: S-0410
tags: [cli]
touches: [flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go]
---
# T-1707 The opencode configuration carries a provider entry derived from the host's providers map, with the key named and never written, and refuses an api OpenCode cannot speak

## Work

- In `opencodeconfig.go`, for a story whose agent names a provider, write `provider.<name>` from the host's entry, with the manifest's overrides (ADR-0129): the package for its `api` (`@ai-sdk/openai-compatible` for `openai-chat`, `@ai-sdk/openai` for `openai-responses`, `@ai-sdk/anthropic` for `anthropic-messages` where S-0406 found it takes a `baseURL`), `options.baseURL`, `options.apiKey` as `{env:<key_env>}`, and `models`. Write the story's model as `<name>/<model>`.
- Refuse an `api` OpenCode cannot speak, naming it and the provider.
- Test each `api`, an override from the manifest, a refused `api`, and that the written file holds the variable's name and never its value.
- Waits for nothing in this story.

## Done when

The provider entry is written as the host's entry says, no key is in it, and the tests pass.

## Notes
