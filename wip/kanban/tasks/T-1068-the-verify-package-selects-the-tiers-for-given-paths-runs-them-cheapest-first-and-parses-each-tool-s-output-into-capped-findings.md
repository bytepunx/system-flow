---
id: T-1068
type: task
nature: improvement
title: The verify package selects the tiers for given paths, runs them cheapest first, and parses each tool's output into capped findings
status: in-progress
parent: S-0273
owner: alex
created: 2026-10-06T22:51:47Z
updated: 2026-10-06T23:52:59Z
transitions:
  - to: ready
    at: 2026-10-06T23:52:59Z
    by: agent-S-0273
  - to: in-progress
    at: 2026-10-06T23:52:59Z
    by: agent-S-0273
stream: S-0273
tags: [cli, testing]
touches: [flai/internal/verify]
---
# T-1068 The verify package selects the tiers for given paths, runs them cheapest first, and parses each tool's output into capped findings

## Work

This is the engine `flai test`, the MCP tool `test`, and the host method `test.run` share. S-0270's `flai verify` builds on it.

- Create the package `flai/internal/verify` with its own `Tier` type: name, command, paths, format, and all-only. The command task maps the manifest's tiers onto it, so this task waits for nothing.
- Selection: given paths or Go packages, choose the tiers whose `paths` globs match, in list order. Fill each command's placeholder with the matching Go packages or files. `--all` takes every tier.
- Running: run each tier in a given directory, a worktree or the main checkout, with its duration. Stop at the first failing tier.
  - Do not shell out to a whole script when the tier is finer: a tier is one command.
  - Model the process handling on `flai/internal/serve/checks.go`'s `runOneCheck`, which kills the whole process group, so a cancelled MCP call or host call leaves nothing running.
- Findings: parse each format into findings of test or check name, path, line, and message.
  - `go-test-json`: from `go test -json`, take the failing test, the `file:line` of its first failure, and the message.
  - `vitest-json`: from vitest's JSON reporter.
  - `golangci-json`: from golangci-lint's `--output.json.path stdout`.
  - `gofmt-list`: from `gofmt -l`, one finding per file.
  - `plain`: the exit status and the last lines of output.
- Cap the findings at a small number, five by default, and count the rest. The answer is never the log.
- Result: overall pass or fail, and per tier its name, state (passed, failed, or not reached), duration, and findings. Make it JSON-ready for the CLI, MCP, and host channel.
- Keep tool outputs as fixtures under `testdata/`. This is why the story keeps a folder touch here.

## Done when

- [ ] Selection picks the right tiers and packages for Go files, flaiover files, and markdown, and `--all` picks every tier (table tests).
- [ ] Each of the five formats turns a recorded failing output into the expected findings, capped with the rest counted (fixture tests).
- [ ] A run stops at the first failing tier, marks later tiers not reached, and records durations. Cancelling the context kills the tier's process group.
- [ ] `go test -race ./internal/verify/...` and golangci-lint pass.

## Notes

Drafted by the planner for S-0273.
