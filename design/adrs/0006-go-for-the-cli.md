---
id: ADR-0006
title: Go for the flai CLI
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0006 Go for the flai CLI

## Context

`flai` must install as a single binary on Linux, macOS, and Windows, run fast, render templates, drive git and docker, and be easy for agents to extend.

## Decision

`flai` is written in Go 1.26 with Cobra for commands, goccy/go-yaml for YAML, and Charm's huh and lipgloss for prompts and tables. Released with GoReleaser. Details in `design/tech/go.md` and `design/tech/go-libraries.md`.

## Consequences

- Metrics have a Go reference implementation that the TypeScript dashboard must match, verified by a fixture test.
- No code sharing with the dashboard. The shared contract is the markdown schema, which is the point.

## Alternatives considered

- Rust: slower iteration, smaller CLI prompt ecosystem.
- TypeScript: shares code with the dashboard but needs a runtime and is harder to distribute as one file.
