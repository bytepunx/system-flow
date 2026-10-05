---
id: S-0273
type: story
nature: improvement
title: flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:32Z
updated: 2026-10-05T01:35:32Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0273 flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings

## Goal

Between edits an agent runs `go test`, `vitest`, `golangci-lint`, or `gofmt` by hand (1,449 turns across 108 runs) and reads the log that comes back, tail and all. `flai test [path|package]...` runs the tiers the project defines for the paths given (the Go packages, the flaiover unit tests, the lint and format checks that apply), cheapest first, and answers pass, or the first failures as findings: test name, path, line, and the assertion's message, with the step's duration; `--all` runs every tier. Over MCP it is `test`, on the host channel `test.run`. `flai verify` (its sibling story) runs it for the whole diff before review.

## Acceptance criteria
- [ ] `flai test` with paths or packages runs the matching tiers from `scripts/` and answers pass or findings as text and `--json`, never more than the first failures and their messages
- [ ] The same is `test` over MCP and `test.run` on the host channel
- [ ] The conventions, the harness prompt, `design/system/devex.md`, `flai-cli.md`, and the user guide send the agent to it for test runs between tasks, and the project's `Makefile` targets keep working for people

## Tasks

## Notes

From the epic's log classification: run tests 797 turns and 96 minutes, lint/format/check 652 turns and 65 minutes.
