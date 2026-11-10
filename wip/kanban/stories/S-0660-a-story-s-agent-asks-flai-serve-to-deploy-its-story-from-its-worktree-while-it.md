---
id: S-0660
type: story
nature: feature
title: A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:08Z
updated: 2026-11-10T10:54:40Z
transitions: []
tags: [cli]
topics: [kluster, deployment]
touches: [flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [S-0659]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 150
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 13759
          cache_read: 1250556
          cache_write: 84563
          cost: 1.202
    - kind: orchestrator
      seconds: 258
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 57
          output: 425
          cache_read: 1932204
          cache_write: 312121
          cost: 0.6087
        - model: claude-sonnet-5-5
          input: 10
          output: 870
          cache_read: 78547
          cache_write: 53641
          cost: 0.0991
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-10T09:34:30Z
  value: 1250
  by: planner-S-0660
  at: 2026-11-10T09:36:53Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T16:54:00Z
  basis: "Its own forecast of 1h30m; 6th in the pull order with an in-progress limit of 3, behind S-0665, S-0708, S-0713, S-0667, S-0674, S-0704, S-0241 and S-0659."
  by: flai
  at: 2026-11-10T10:54:40Z
finalized:
  by: orchestrator
  at: 2026-11-10T09:39:54Z
---
# S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled

## Goal

A story's agent asks flai serve to deploy its story from its worktree while the story is in progress, and the namespace stays until the story is accepted or cancelled. The agent has the running system to work against, and runs no kluster command that writes.

## Acceptance criteria

- [ ] With the `kluster` host action on and a `deploy` set, a story's agent asks for a deploy of its own story in progress with one flai command and one MCP tool, and flai serve builds in the story's worktree at its branch's head, loads, and deploys into `s-nnnn`, as it does for a story in review
- [ ] The request waits for the job and answers with its outcome, the commit, each step's result, and the log's path; a request that outlasts its wait says the job is still queued or running and how to read it later
- [ ] A request is refused, with the reason and what to do, for another story, a story not in progress or review, uncommitted changes in the worktree, the action off, no `deploy`, and a request over the limit on story namespaces
- [ ] `flai guard` still refuses every agent `flai kluster up`, `down`, `deploy`, `teardown`, and `promote`, and refuses a sub-agent the new request
- [ ] A story in progress keeps its namespace: flai serve tears it down only when the story is accepted or cancelled, or when its agent or the operator asks
- [ ] A story that enters review at the head already deployed is not deployed again, and one at a new head is deployed in place
- [ ] A failed deploy asked for by the agent opens no thread to the operator: the agent gets the step, the exit code, and the log
- [ ] Each job the agent asked for is in the host journal and the namespace's record with who asked
- [ ] `design/system/kluster.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` and its reference say all of it

## Tasks
- T-2876 A request for a story's deploy or teardown is a file flai serve takes and answers, and the namespace's record names who asked
- T-2879 The kluster block sets how many story namespaces may live at once
- T-2882 flai guard refuses a sub-agent the request and still refuses every agent kluster up, down, deploy, teardown, and promote
- T-2884 flai serve takes a story's request, deploys a story in progress at its head, and keeps its namespace until acceptance, cancellation, or a teardown asked for
- T-2887 A flai kluster subcommand and an MCP tool ask flai serve for the session's own story's deploy or teardown and print the answer
- T-2889 kluster.md, flai-cli.md, and the user guide say how a story's agent asks for its deploy and when its namespace goes

## Notes

- `flai/internal/serve/kluster.go` is the pass, and `flai/internal/guard/guard_rules.go` holds `klusterWrites`.

### Planning

`after`: S-0659. This story builds the design S-0659's ADR decides: the request's command and MCP tool names, how the agent waits behind the job queue and for how long, where the project sets the limit on story namespaces and its default, and the table of when a namespace is made, deployed again, and torn down that replaces ADR-0254 §5's. The tasks name these as the ADR will. Where it decides otherwise, the story's agent follows the ADR and changes the tasks' words.

Touches, all files, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/kluster.go`, `flai/internal/guard/guard_rules.go` | declared (Notes) | The pass, and `klusterWrites` |
| `design/system/kluster.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` | design (criterion 9) | Named by the criterion; the reference is generated by `make flai-reference` |
| `docs/operators/settings.md`, `design/system/project-manifest.md` | co-change (`flai touches suggest`, 18% for settings.md) and design | The limit on story namespaces is a manifest key, which both list; the flag index in settings.md is generated with the reference |
| `flai/internal/kluster/request.go`, `request_test.go` | layout | A new file for the request, modelled on `flai/internal/closeout/closeout.go`'s `Ask` with its take and answer waits, the hand-off `flai verify --job` makes |
| `flai/internal/kluster/deploy.go`, `deploy_test.go` | layout | `Record` and `Job` gain who asked |
| `flai/internal/manifest/manifest_kluster.go`, `settings.go`, `manifest_test.go`, `settings_test.go` | layout | The limit key in the `kluster:` block |
| `flai/internal/serve/kluster_test.go` | layout | The pass's tests |
| `flai/cmd/kluster.go`, `kluster_test.go` | layout | The new subcommand beside `deploy` and `teardown` |
| `flai/internal/mcpserver/kluster.go`, `kluster_test.go`, `folder.go`, `server_test.go` | layout | The new MCP tool, registered in `folder.go`, and `TestToolsAreAdvertised`'s list |
| `flai/internal/guard/kluster_test.go` | layout | `TestFlaiKlusterWritesAreTheOperators` gains the request |

Left out on purpose: `flai/internal/hostapi/kluster.go`, since the request goes through files under `.flai-cache` as the close-out's does, not the host API; and `flai/internal/guard/guard_tables.go`, since a sub-agent is refused a tool and a subcommand that no read table lists. `docs/operators/index.md` (19% co-change on the second run) is left to S-0663's operator side; nothing else `flai touches suggest` lists reaches 15%.

Forecast: 1h30m, raised by half from the 1h `flai forecast` gives for size 33, because S-0577 to S-0580, the kluster stories before it, each took 39m to 61m with 6 or 7 criteria, and this one has 9 across six parts of flai. The delivery, 2026-11-10T15:25:00Z, is flai's from that duration and the pull order.

Cost of delay: 1250 USD a week, as `flai cod S-0660` works it out from the operator's `revenue_per_week: 1250` on this story (TH-0979, after the orchestrator's recommendation there, from D-0479). Not adjusted: the operator set the figure on this story directly, and this story is the one that lets a kluster story's agent deploy before review, which is what holds up the signet and authstar releases that revenue waits on.
