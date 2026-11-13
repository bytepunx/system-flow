---
id: T-2882
type: task
nature: feature
title: flai guard refuses a sub-agent the request and still refuses every agent kluster up, down, deploy, teardown, and promote
status: done
parent: S-0660
owner: alex
created: 2026-11-10T08:27:24Z
updated: 2026-11-13T08:02:50Z
transitions:
  - to: ready
    at: 2026-11-13T07:33:53Z
    by: agent-S-0660
  - to: in-progress
    at: 2026-11-13T07:33:53Z
    by: agent-S-0660
  - to: done
    at: 2026-11-13T08:02:50Z
    by: agent-S-0660
stream: S-0660
tags: [cli]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go]
usage:
  source: log
  seconds: 1737
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 27895
      cache_read: 4555149
      cache_write: 142890
      cost: 2.2946
---
# T-2882 flai guard refuses a sub-agent the request and still refuses every agent kluster up, down, deploy, teardown, and promote

## Work

ADR-0292 §2 and criterion 4. In `flai/internal/guard/guard_rules.go`:

- `flai kluster deploy S-nnnn --job` and `flai kluster teardown S-nnnn --job` are allowed to a story's agent, served or not, for its own story only (`FLAI_STORY`, else the story in `agent-S-nnnn`), and refused to every sub-agent and for any other story or none, with the reason and what to do.
- `deploy` and `teardown` without `--job`, and `up`, `down`, and `promote` with any flags, stay refused to every agent, as `klusterWrites` refuses them now. Check that no flag form, such as `--job=false`, a flag before the subcommand, or `--json` alone, reaches a write.
- The MCP tools `kluster_deploy` and `kluster_teardown` are refused to every sub-agent, as `verify` is (find how the guard refuses it), and to the story's agent for another story when the tool names one.

A sub-agent is refused a tool and a subcommand no read table lists; if `guard_tables.go` needs a change, widen this task's touches.

It waits for nothing: the names come from ADR-0292, so it runs in the first layer.

## Done when

- `flai/internal/guard/kluster_test.go` shows the story's agent, served or not, allowed the request for its own story by command and by tool; refused it for another story; a sub-agent refused both; and every agent still refused `up`, `down`, `promote`, and `deploy` and `teardown` without `--job`.
- `flai test` on the changed paths passes.

## Notes

- Out of scope, as agreed on TH-0987: a sub-agent's `flai kluster status` falls through to the generic refusal in `flai/internal/guard/guard.go`, against `design/system/kluster.md`'s guard section. The story's agent records it as an issue; do not fix it here.
