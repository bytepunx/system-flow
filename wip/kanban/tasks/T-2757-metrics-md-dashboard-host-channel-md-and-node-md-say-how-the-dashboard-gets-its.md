---
id: T-2757
type: task
nature: improvement
title: metrics.md, dashboard-host-channel.md, and node.md say how the dashboard gets its metrics, that the host channel is built, and which lockfile the build installs from
status: done
parent: S-0626
owner: alex
created: 2026-11-09T17:13:16Z
updated: 2026-11-09T17:33:43Z
transitions:
  - to: ready
    at: 2026-11-09T17:27:15Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:27:16Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:33:43Z
    by: agent-S-0626
stream: S-0626
tags: [design, docs]
touches: [design/system/metrics.md, design/system/dashboard-host-channel.md, design/tech/node.md]
usage:
  source: log
  seconds: 387
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 29
      output: 209
      cache_read: 2318403
      cache_write: 26392
      cost: 0.856
---
# T-2757 metrics.md, dashboard-host-channel.md, and node.md say how the dashboard gets its metrics, that the host channel is built, and which lockfile the build installs from

## Work

- `design/system/metrics.md` § Data access, line 605: the dashboard scans no `wip/kanban`. It asks flai on the host for items and stats and keeps each answer until flai says a file changed (`flaiover/src/lib/server/stats.ts:1-5`, `repo.ts:916-921`; ADR-0016, ADR-0029). Its fixture test runs against `flai/internal/metrics/testdata/good`, not a template sample (`stats.test.ts:12`). Correct both sentences. This section says how the data is read and defines no metric, so it needs no ADR: say so in the narrative's `## Decisions`. S-0624 also touches this document, in its Close-outs section.
- `design/system/dashboard-host-channel.md` line 10, "Nothing here is built", and line 149, "S-0082, in progress": the channel is built (`flai/internal/hostapi/`, `flaiover-dashboard.md` § The channel to flai on the host), and S-0082 is done and archived. Say what landed, with its stories.
- `design/tech/node.md` line 20, "a single lockfile format": the build and CI install from `flaiover/pnpm-lock.yaml` alone (`flaiover/Dockerfile:20`, `.github/workflows/flaiover.yml:33`, `system-flow.yaml`). Say that. `flaiover/package-lock.json` and the root `package-lock.json` are still tracked, and nothing reads them. S-0630 removes them (I-0297), so this task does not.

Bump each document's `updated`.

It waits for no task: no other task changes these files.

## Done when

- `metrics.md` names no scan of `wip/kanban` and no template sample repository.
- `dashboard-host-channel.md` names S-0082 as done and no longer says nothing is built.
- `node.md` names the lockfile the build installs from.
- `flai check --strict` and the markdown lint pass on the three files.

## Notes
