---
id: T-3621
type: task
nature: improvement
title: The MCP tool accept takes keep_removals, each a path and a reason, and passes them to flai accept --by orchestrator on the host
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:57:52Z
updated: 2026-11-27T23:57:52Z
transitions: []
stream: S-0852
tags: [flai]
touches: [flai/internal/mcpserver/accept.go, flai/internal/mcpserver/accept_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/orchestrate_test.go]
after: [T-3620]
---
# T-3621 The MCP tool accept takes keep_removals, each a path and a reason, and passes them to flai accept --by orchestrator on the host

## Work

In I-0415 the orchestrator's `accept` tool had no way to give a keep (ADR-0278). Add `keep_removals` to its input in `flai/internal/mcpserver/accept.go`, each a path and a reason, as `accept.run` takes them for the operator. Pass each as a `--keep-removal` and `--reason` pair to the `flai accept <S-nnnn> --by orchestrator` the tool runs on the host, with `dry_run` too, so the preview shows the paths kept. flai's refusal for a path the narrative does not name, or with the permission off, comes back as the tool's error. Update the tool's description to name the field, the permission, and the narrative condition.

In `flai/internal/guard/guard_rules.go`, update the orchestrator's accept rule text, which says an undone path waits for the operator. The guard still refuses the tool to every role but the orchestrator.

It waits for T-3620: the tool passes the flags that task makes flai accept take.

## Done when

- A test in `flai/internal/mcpserver/accept_test.go` shows `keep_removals` reaching `flai accept` as flag pairs, in order, and with `dry_run`.
- A test shows flai's refusal for a keep with the permission off answered as the tool's error.
- `flai test` passes on the changed paths.

## Notes
