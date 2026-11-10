---
id: TH-1003
title: "Plan for S-0665: record a container run's image, container, and mounts in agents.json and name them in the journal; five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0665-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md
  item: S-0665
status: open
participants: [planner-S-0665]
created: 2026-11-10T08:51:49Z
updated: 2026-11-10T08:51:49Z
---

# TH-1003 Plan for S-0665: record a container run's image, container, and mounts in agents.json and name them in the journal; five tasks in four layers

On wip/kanban/stories/S-0665-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md.

## Entries

### 2026-11-10T08:51:49Z planner-S-0665
**Proposed solution.** Today the container adapter builds the `docker run` argv, with its image, `--name`, and `-v` mounts, then drops them. `agents.json` keeps only `command: docker`. The fix:

1. `harness.Start` carries the image, the container's name, and the mounts.
2. `startOf`, which every story and strategic run goes through, copies them onto the run as `image`, `container`, and `mounts`.
3. The journal's start line names the image and the container.

**Tasks and layers.**

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2942 The container adapter's Start names the image, name, and mounts | `flai/internal/harness/harness.go`, `container.go`, `container_test.go` | none |
| 2 | T-2943 flai serve records them in agents.json, with a test reproducing I-0274 | `flai/internal/serve/agents.go`, `agents_start.go`, `container.go`, `container_test.go` | T-2942 |
| 3 | T-2945 The journal's start line of a container run names its image and container | `agents_start.go`, `plan.go`, `orchestrate.go`, `analyze.go`, `curate.go`, `container_test.go` under `flai/internal/serve/` | T-2943 |
| 4 | T-2946 Design and user guide | `design/system/flai-cli.md`, `docs/users/flai.md`, `design/system/strategic-agents.md` | T-2943, T-2945 |
| 4 | T-2948 Close I-0274 | I-0274's file, `design/issues/summary.md` | T-2943, T-2945 |

T-2945 waits for T-2943 because it reads the run's new fields and shares `agents_start.go` and `container_test.go` with it. T-2946 and T-2948 share no path and run together.

**Figures.**

- Forecast: 1h30m, delivery 2026-11-10T13:05:00Z. `flai forecast` gave 31m; I raised it for the five start paths, the fake-docker test, and the close-out's integration and smoke tiers.
- Cost of delay: 37.50 USD a week, from flai's input of 15m per cycle. It stands as `flai cod` gives it.

**Assumptions.**

- The mounts go in `agents.json` only, not the journal line, which would grow too long to read.
- No dashboard page shows the new fields. `agent.status` already returns the whole run, so they reach the wire; a page for them would be a story of its own.
- `flai/internal/serve/role.go` is untouched, since a role a project defines is refused on the container harness.
- The story has no `after`: it reads no other open story's output. S-0670, S-0682, and S-0687 work on containers too, but on `scripts/agent-container.sh`, pilot, and the strategic container's configuration, not this record.
