---
id: S-0273
type: story
nature: improvement
title: flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings
status: in-progress
parent: E-0017
owner: alex
created: 2026-10-05T01:35:32Z
updated: 2026-10-06T23:52:58Z
transitions:
  - to: ready
    at: 2026-10-06T23:32:45Z
    by: alex
  - to: in-progress
    at: 2026-10-06T23:51:15Z
    by: agent-S-0273
tags: [cli, mcp]
topics: [automation, mcp, hostapi, code, conventions, template]
touches: [flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/root.go, flai/internal/verify, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, system-flow.yaml, template/root/system-flow.yaml.tmpl, scripts/flai-test.sh, scripts/test.sh, scripts/flaiover-unit.sh, scripts/with-env.sh, scripts/README.md, Makefile, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/conventions/work-management.md, design/conventions/tooling.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/tooling.md, template/CHANGELOG.md, design/system/devex.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/dashboard-host-channel.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 213
  by: planner-E-0017
  at: 2026-10-06T11:36:16Z
forecast:
  duration: 75m
  delivery: 2026-10-07T01:05:00Z
  basis: "Its own forecast of 1h15m; 1st in the pull order with an in-progress limit of 3, behind S-0303."
  by: flai
  at: 2026-10-06T23:47:21Z
finalized:
  by: alex
  at: 2026-10-06T22:48:41Z
---
# S-0273 flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings

## Goal

Between edits an agent runs `go test`, `vitest`, `golangci-lint`, or `gofmt` by hand (1,449 turns across 108 runs) and reads the log that comes back, tail and all. `flai test [path|package]...` runs the tiers the project defines for the paths given (the Go packages, the flaiover unit tests, the lint and format checks that apply), cheapest first, and answers pass, or the first failures as findings: test name, path, line, and the assertion's message, with the step's duration; `--all` runs every tier. Over MCP it is `test`, on the host channel `test.run`. `flai verify` (its sibling story) runs it for the whole diff before review.

## Acceptance criteria
- [ ] `flai test` with paths or packages runs the matching tiers from `scripts/` and answers pass or findings as text and `--json`, never more than the first failures and their messages
- [ ] The same is `test` over MCP and `test.run` on the host channel
- [ ] The conventions, the harness prompt, `design/system/devex.md`, `flai-cli.md`, and the user guide send the agent to it for test runs between tasks, and the project's `Makefile` targets keep working for people

## Tasks
- T-1067 The manifest declares a project's test tiers: name, command, the paths that select each, and the format of its output
- T-1068 The verify package selects the tiers for given paths, runs them cheapest first, and parses each tool's output into capped findings
- T-1069 flai test runs the manifest's tiers for paths or packages, or every tier with --all, and prints pass or the first findings as text and --json
- T-1070 The MCP tool test runs the tiers for paths in a story's worktree and answers pass or the first findings
- T-1073 The host method test.run runs flai test in a story's worktree behind the checks host action and answers its result
- T-1079 The conventions, the harness prompt, the design, and the user guide send the agent to flai test for test runs between tasks

## Notes

From the epic's log classification: run tests 797 turns and 96 minutes, lint/format/check 652 turns and 65 minutes.

### Planning

Planned by planner-S-0273 on 2026-10-06, revisiting planner-E-0017's enrichment. The plan is on this story's plan thread.

Tasks, in four layers:

1. T-1067, the manifest's `tests` key, and T-1068, the `flai/internal/verify` runner and parsers. They share no path.
2. T-1069, the command `flai test`, and T-1070, the MCP tool `test`. Both wait for layer 1 and share no path.
3. T-1073, the host method `test.run`. It waits for T-1069, whose `flai test --json` it builds.
4. T-1079, the conventions, the prompt, and the docs. It waits for T-1069, T-1070, and T-1073, which it describes, and for T-1067, with which it shares `template/CHANGELOG.md`.

Touches. `flai touches suggest S-0273` co-changes nothing above 18%: the dashboard design and docs, the ADR index, and `mcpserver/server.go`. None of them is a path this story changes, so none is added. Every earlier touch is kept, with these changes:

- `flai/internal/manifest` (a folder) is narrowed to `manifest.go`, `manifest_test.go`, `settings.go`, and `settings_test.go` (layout). The tiers are a field beside `Checks` in `manifest.go`, and `flai manifest set` writes them through the catalog in `settings.go`.
- Added `system-flow.yaml` (layout): this repository declares its tiers there, through `flai manifest set`.
- Added `template/root/system-flow.yaml.tmpl` (design): a project made from the template declares its tiers there.
- Added `flai/internal/mcpserver/test.go` and `test_test.go` (layout): each tool has its own file, as `criteria.go` has.
- Added `design/system/dashboard-host-channel.md` (design): it lists the host methods, `checks.run` among them.
- Added `design/conventions/tooling.md` and its `template/root` copy (design, criterion 3): line 43 lists the `make` test targets.
- Folder touch kept: `flai/internal/verify` (layout). It is a new package, and its files and `testdata/` fixtures cannot be named until T-1068 writes them. S-0270 extends it after this story, and no other open story touches it.
- Kept from before:
  - `flai/cmd/test.go`, `test_test.go`, and `root.go` (layout).
  - `scripts/flai-test.sh`, `scripts/test.sh`, `scripts/README.md`, and `Makefile` (declared seeds and co-change). The scripts get a comment naming the tiers they mirror. The `Makefile` is checked and changed only if a target breaks.
  - `flai/internal/mcpserver/folder.go`, `flai/internal/hostapi/hostapi.go`, `writes.go`, and `writes_test.go` (layout). `test.run` sits beside `checks.run` in `writes.go`.
  - `flai/internal/harness/harness.go` and `harness_test.go`, `delegation.md`, `work-management.md`, their `template/root` copies, and `template/CHANGELOG.md` (design, criterion 3).
  - `design/system/devex.md`, `project-manifest.md`, `flai-cli.md`, `docs/operators/settings.md`, `docs/users/flai.md`, and `flai-reference.md` (design).
- Left out:
  - `scripts/env.sh`, `scripts/install-tools.sh`, the CI workflows, and the old `design/issues` files that co-changed. They are about installing tools, not running tiers.
  - `flai/internal/serve/checks.go`. T-1068 models its process handling on `checks.go` but does not change it. S-0270 decides whether the two runners merge.

Topics: added `conventions` and `template`. T-1079 changes three conventions and their template copies.

Forecast: 75m, against flai's 40m (83 s per unit over 25 done large improvement stories, times size 29) and planner-E-0017's 60m. The size counts touches, but the work is in six tasks across four layers:

- a manifest key with validation and `flai manifest set` support;
- five output formats to parse into findings;
- three surfaces: the CLI, MCP, and the host channel.

S-0217 took 68m for three commands that parse nothing. The delivery is flai's, 2026-10-07T08:12Z, moved by the added 35m to 08:47Z.

Cost of delay: 213 USD a week stands, against flai's 141.36. flai shares E-0017's 900 USD a week (alex's 6h a cycle, TH-0171) by forecast duration. The value comes from the hand turns each story removes, so the share by turns is truer. This story removes 1,449 of about 6,100 hand turns, the second most in the epic: 1,449 / 6,100 × 900 ≈ 213.
