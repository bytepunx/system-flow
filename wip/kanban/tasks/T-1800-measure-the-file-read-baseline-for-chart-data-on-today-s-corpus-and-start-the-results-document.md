---
id: T-1800
type: task
nature: experiment
title: Measure the file-read baseline for chart data on today's corpus and start the results document
status: in-progress
parent: S-0427
owner: alex
created: 2026-10-16T05:59:07Z
updated: 2026-10-16T06:00:57Z
transitions:
  - to: ready
    at: 2026-10-16T06:00:57Z
    by: agent-S-0427
  - to: in-progress
    at: 2026-10-16T06:00:57Z
    by: agent-S-0427
stream: S-0427
tags: [cli]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
---
# T-1800 Measure the file-read baseline for chart data on today's corpus and start the results document

## Work

Start `design/experiments/S-0427-sqlite-and-postgres-read-models.md` from `design/experiments/template.md`, with its front matter (`title`, `updated`, `status`, `story: S-0427`). Write `## Hypothesis` and `## Success measure`: a read model loaded from the files answers the chart data's broad reads faster than flai's file reads, by a margin worth its storage, RAM, and CPU, for both the operator's wait and an agent's one-shot call.

Then measure the baseline flai gives today, on the main checkout's corpus, and record it under `## What was done`:

- The corpus: the count of markdown files under `wip/kanban`, `wip/archive`, `wip/threads`, and `wip/messages`, and their total size.
- `stats.get`, the dashboard's chart data (`flai/internal/hostapi/reads.go`, through `statsread.Read` and `metrics.Compute`): `flai hostapi --timing stats.get '<params>'`, cold and warm, the median of three, as `design/system/server-performance.md` measures it, with its `repo.list` and `stats.compute` phases.
- The one-shot call an agent makes: `flai stats --json`, wall time, the median of three.
- Max RSS and CPU time of each, with `/usr/bin/time -v`.

Save one `flai stats --json` output in a scratch folder the repository ignores, such as `bin/experiments/S-0427/`, as the reference the later tasks check their queries against. Commit no code and no data.

Waits for no task: it is the first layer.

## Done when

- The results document exists with its front matter, `## Hypothesis`, and `## Success measure`.
- `## What was done` gives the corpus size and the baseline's cold and warm timings, phases, RSS, and CPU, each with the command that measured it.
- The reference `flai stats --json` output is saved in the scratch folder.

## Notes
