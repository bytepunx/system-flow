---
id: T-2418
type: task
nature: feature
title: flai serve settings prints the host settings the dashboard's settings.get gives, its mcp block naming the bind address flai serve mcp sets
status: backlog
parent: S-0506
owner: alex
created: 2026-10-31T07:58:50Z
updated: 2026-10-31T07:58:50Z
transitions: []
stream: S-0506
tags: [cli, mcp]
touches: [flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/cmd/serve.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md]
after: [T-2152]
---
# T-2418 flai serve settings prints the host settings the dashboard's settings.get gives, its mcp block naming the bind address flai serve mcp sets

## Work

Criterion 3 says `flai serve settings` shows the bind address, and no such command exists: each `flai serve` area shows only its own settings. Add `flai serve settings` in a new `flai/cmd/serve_settings.go`, registered beside the other subcommands in `flai/cmd/serve.go`. It prints what `hostSettings` in `flai/cmd/serve_actions.go` gives the dashboard's `settings.get` for the project in the current folder: as JSON with `--json`, and as a readable summary without it. Add the bind address T-2152 stores to the `mcp` block of `hostSettings`, beside `running`, `url`, and `pid`, as set or as the loopback default. The command prints no token, as `hostSettings` gives none. Add a line on the command to the `flai serve` section of `docs/users/flai.md` and to `docs/operators/index.md`, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

`flai serve settings` prints the bind address `flai serve mcp bind` set, and the loopback default when it is unset. A test in `flai/cmd/serve_settings_test.go` checks both and that no token appears. `TestReferenceIsCurrent` passes, and so does `flai test flai/cmd`, `settings_doc_test` aside, since T-2167 gives the key its row.

## Notes

Added by planner-S-0506 on 2026-10-31 after the operator's edit of criterion 3 kept `flai serve settings` as what shows the setting.
