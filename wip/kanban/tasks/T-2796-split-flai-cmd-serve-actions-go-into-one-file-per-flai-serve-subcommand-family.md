---
id: T-2796
type: task
nature: improvement
title: Split flai/cmd/serve_actions.go into one file per flai serve subcommand family
status: backlog
parent: S-0631
owner: alex
created: 2026-11-09T22:04:40Z
updated: 2026-11-09T22:04:40Z
transitions: []
stream: S-0631
tags: [cli]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_agent.go, flai/cmd/serve_checks.go, flai/cmd/serve_components.go, flai/cmd/serve_harness.go, flai/cmd/serve_provider.go, flai/cmd/serve_notify.go, flai/cmd/serve_mcp.go, flai/cmd/serve_view.go, flai/cmd/serve_agent_run.go]
---
# T-2796 Split flai/cmd/serve_actions.go into one file per flai serve subcommand family

## Work

`serve_actions.go` is 2,754 lines. Move each `flai serve` subcommand family into its own file of package `cmd`. What any of them does stays the same. The line ranges below are approximate, from the plan's read on 2026-11-09:

| Lines | Declarations | File |
|-------|--------------|------|
| 42–578 | `app.host`, `hostBudget`, `journal`, `hostActions`, `servedRoots`, `newServeEnableCmd`, `newServeActionsCmd`, `newServeJournalCmd`, the kluster scaffold (`klusterScaffold`, `scaffoldKluster`, `commitScaffold`), and `serveHost` (about 1186) | `serve_actions.go`, which keeps them |
| 579–1021 | `newServeAgentCmd`, `agentScope`, `showAgentCommand`, `printLifts`, `agentConfig` | `serve_agent.go` |
| 1022–2013 | the checks, components, harness, provider and forge, and notify subcommands | `serve_checks.go`, `serve_components.go`, `serve_harness.go`, `serve_provider.go`, `serve_notify.go` |
| 2014–2349 | the MCP subcommands, and the served view | the existing `serve_mcp.go` (104 lines now), and `serve_view.go` |
| 2350–2754 | restart, start, commit, stop, orchestrate, stream, `mcpAgents`, `agentNow` | `serve_agent_run.go` |

`serve_settings.go`'s `showHostSettings` is a different view from `hostSettings` (about 2181): leave `serve_settings.go` alone. The new names do not collide with any existing `flai/cmd/serve_*.go`.

Change no test file. `serve_actions_test.go` stays as it is.

This task waits for nothing and shares no path with the other splits, so all of them run together.

## Done when

- No non-test file in `flai/cmd` is over 2,000 lines, and none of the new files is over 1,200.
- `flai serve --help` and each moved subcommand's `--help` print what they printed before.
- `git diff --stat` shows no `_test.go` file changed.
- `flai test flai/cmd/serve_actions.go` passes, gofmt, goimports, and golangci-lint included.

## Notes
