---
id: T-2882
type: task
nature: feature
title: flai guard refuses a sub-agent the request and still refuses every agent kluster up, down, deploy, teardown, and promote
status: backlog
parent: S-0660
owner: alex
created: 2026-11-10T08:27:24Z
updated: 2026-11-10T08:32:02Z
transitions: []
stream: S-0660
tags: [cli]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go]
---
# T-2882 flai guard refuses a sub-agent the request and still refuses every agent kluster up, down, deploy, teardown, and promote

## Work

Make `flai guard` let a story's agent make the request, by the subcommand and MCP tool S-0659's ADR names, and refuse it to a sub-agent. Keep `klusterWrites` in `flai/internal/guard/guard_rules.go` refusing every agent `flai kluster up`, `down`, `deploy`, `teardown`, and `promote`, with any flags. The `klustering` rule matches on the subcommand, so a new subcommand name is not refused to the story's agent; check that no flag form of the request reaches `deploy` or `teardown`. A sub-agent is refused a tool and a subcommand that no read table lists, so `guard_tables.go` should not need a change. If it does, widen this task's touches.

It waits for nothing: the names come from the ADR, not from the code of the CLI task, so it runs in the first layer.

## Done when

- `flai/internal/guard/kluster_test.go` shows the story's agent, served or not, allowed the request by command and by tool; a sub-agent refused both; and every agent still refused `up`, `down`, `deploy`, `teardown`, and `promote`.
- `go test ./internal/guard/...` passes.

## Notes

- Out of scope, as agreed on TH-0987: a sub-agent's `flai kluster status` falls through to the generic refusal in `flai/internal/guard/guard.go`, against `design/system/kluster.md`'s guard section. Record it with `flai issue new` when you start, so it gets its own story; do not fix it here.
