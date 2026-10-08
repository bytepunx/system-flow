---
id: T-1380
type: task
nature: feature
title: The user, operator, and design documents say what a provider sets at start, what it refuses, and where its key's variable must be
status: backlog
parent: S-0350
owner: alex
created: 2026-10-08T08:44:15Z
updated: 2026-10-08T08:44:15Z
transitions: []
stream: S-0350
tags: [cli]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md]
after: [T-1378]
---
# T-1380 The user, operator, and design documents say what a provider sets at start, what it refuses, and where its key's variable must be

## Work

- `docs/users/flai.md`, "Starting an agent": the variables the `claude-code` adapter sets from a provider, the `command` adapter's `FLAI_PROVIDER_*` and `{provider}`, and each refusal with what fixes it.
- `docs/operators/index.md`, under authentication: the key lives in the variable `key_env` names, in the environment `flai host` or `flai serve` runs in, never in a file flai writes; how to set it for a host started at login.
- `design/system/flai-cli.md`: the start through a provider, `EnvFrom`, and why no key reaches a log.

It waits for T-1378, so that it describes the start as built.

## Done when

- The three documents match the code, and no example holds a key.
- `flai test docs/ design/system/flai-cli.md` passes.

## Notes

Layer 3 of S-0350.
