---
id: S-0273
type: story
nature: improvement
title: flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:32Z
updated: 2026-10-06T11:36:34Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, hostapi, code]
touches: [flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/root.go, flai/internal/verify, flai/internal/manifest, scripts/flai-test.sh, scripts/test.sh, scripts/README.md, Makefile, flai/internal/mcpserver/folder.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/devex.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 213
  by: planner-E-0017
  at: 2026-10-06T11:36:16Z
forecast:
  duration: 60m
  delivery: 2026-10-07T02:01:00Z
  basis: "flai's 43m (size 29) raised to 60m because findings are parsed from four tools' outputs (go test, vitest, golangci-lint, gofmt) and tiers become project data, more than S-0217's 68m per command; 35th in the pull order."
  by: planner-E-0017
  at: 2026-10-06T11:34:00Z
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

### Planning

Touches, none declared before. `flai touches suggest S-0273` was seeded with `scripts/flai-test.sh` and `scripts/test.sh`, which 7 commits changed:

- `flai/cmd/test.go`, `flai/cmd/test_test.go`, `flai/cmd/root.go`: layout. A new command, registered in `root.go`.
- `flai/internal/verify`: layout. A new package that runs a tier and parses its output into findings. S-0270 builds on it, so the name is a prediction both stories share.
- `flai/internal/manifest`, `design/system/project-manifest.md`, `docs/operators/settings.md`: design. Today the tiers exist only as scripts and `Makefile` targets. "The tiers the project defines" needs somewhere a project made from the template declares them. This is an assumption on the plan thread.
- `scripts/flai-test.sh`, `scripts/test.sh`, `scripts/README.md`, `Makefile`: declared seeds and co-change (README 3 of 7, Makefile 4 of 7). Machine-readable output, such as `go test -json`, may need flags here, and criterion 3 keeps the targets working.
- `flai/internal/mcpserver/folder.go`, `flai/internal/hostapi/hostapi.go`, `writes.go`, `writes_test.go`: layout. The `test` tool and the `test.run` method.
- `flai/internal/harness/harness.go`, `harness_test.go`, `design/conventions/delegation.md`, `work-management.md`, their `template/root` copies, and `template/CHANGELOG.md`: design (criterion 3). The prompt and the delegation rule name the hand runs.
- `design/system/devex.md`, `flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`: design (criterion 3).
- Left out: `scripts/env.sh`, `scripts/install-tools.sh`, the CI workflows, and the old `design/issues` files that co-changed (2 to 3 of 7). They are about installing tools, not running tiers.

Forecast: flai gave 43m (89 s per unit over 21 done large improvement stories, times size 29). I raised it to 60m because findings come from four tools' outputs: go test, vitest, golangci-lint, and gofmt. The tiers also become project data. S-0217 took 68m for three commands that parse nothing. The delivery is flai's, 2026-10-07T01:44Z, moved by the added 17m.

Cost of delay: 213 USD a week, against flai's 157.89. This is E-0017's 900 USD a week shared by the turns each story removes. This one removes 1,449 hand runs of about 6,100 turns, the second most.
