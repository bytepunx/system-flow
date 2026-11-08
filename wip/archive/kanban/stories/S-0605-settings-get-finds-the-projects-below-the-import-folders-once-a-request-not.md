---
id: S-0605
type: story
nature: improvement
title: settings.get finds the projects below the import folders once a request, not once a served project, so it answers in under 100 ms, not 825 ms
status: done
parent: E-0033
owner: alex
created: 2026-11-07T22:40:46Z
updated: 2026-11-08T01:13:20Z
transitions:
  - to: ready
    at: 2026-11-08T00:11:13Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-08T00:11:31Z
    by: agent-S-0605
  - to: review
    at: 2026-11-08T01:10:32Z
    by: agent-S-0605
  - to: done
    at: 2026-11-08T01:13:20Z
    by: orchestrator
tags: [cli]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_project.go, flai/internal/serve/below.go, flai/cmd/serve_project_test.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/system/server-performance.md, flai/internal/serve/below_test.go, flai/cmd/serve_settings.go]
after: [S-0596]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3554
  priced_by: harness
  turns:
    - day: 2026-11-08
      ceremony: 1
      hand_edits: 3
      work: 68
  models:
    - model: claude-opus-5-5
      input: 202
      output: 51814
      cache_read: 27832124
      cache_write: 671431
      cost: 10.9871
  strategic:
    - kind: orchestrator
      seconds: 182
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 159
          cache_read: 1243952
          cache_write: 126930
          cost: 0.3495
cost_of_delay:
  inputs:
    penalty_per_week: 20
    by: orchestrator
    at: 2026-11-07T23:50:59Z
  value: 20
  by: planner-E-0033
  at: 2026-11-08T00:01:11Z
forecast:
  duration: 30m
  delivery: 2026-11-08T00:33:00Z
  basis: "flai's 8m (79 s per unit of size over 20 small improvement stories, size 6), raised to 30m: two code changes, one threading a context through the host API, and a measurement on the operator's host, as S-0606's comparable stories took 12 to 36 min."
  by: planner-S-0605
  at: 2026-11-07T23:47:52Z
finalized:
  by: alex
  at: 2026-11-07T23:49:40Z
---
# S-0605 settings.get finds the projects below the import folders once a request, not once a served project, so it answers in under 100 ms, not 825 ms

## Goal

Cause 15 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `settings.get` has 4,104 slow calls in 4 days at p50 825 ms. `servedView` (`flai/cmd/serve_actions.go`) asks `importRootOf` (`flai/cmd/serve_project.go`) for each of the 13 served projects, and each runs `serve.FindBelow` (`flai/internal/serve/below.go`), 45 to 85 ms: 13 walks of the import folders and 169 manifest loads a request. It takes no context, so it records no phases.

## Acceptance criteria

- [x] `settings.get` works out the projects below the import folders once a request, or takes what the running flai serve already found, and answers what it answered before, which a test with three served projects under two import folders shows.
- [x] `settings.get` records its phases.
- [x] `settings.get` on the operator's host answers under 100 ms, and server-performance.md records the before and after.

## Tasks
- T-2639 servedView works out the projects below the import folders once a request, from the running flai serve's status when it runs, and a test with three served projects under two import folders shows the same answer
- T-2645 settings.get passes its context to Host.Settings and records its phases: the configuration, the served projects, and the project's own settings
- T-2654 Measure settings.get before and after on the operator's host, and record cause 15's fix and figures in server-performance.md

## Notes

### Planning

Planned by planner-S-0605 on 2026-11-07. No thread on S-0605 asked for changes; the plan follows the story's words and cause 15 in `design/system/server-performance.md`.

What the code shows: `hostSettings` (`flai/cmd/serve_actions.go`) calls `servedView`, which calls `listServedProjects` once and then `servedBelow` for each served project. `servedBelow` calls `importRootOf`, which loads the host configuration through `importRoots` and runs `serve.FindBelow` over every import folder, only to ask whether one root is among the found ones. `listServedProjects` itself runs `serve.FindCandidates` over the same folders, and `belowImportRoots` runs `FindBelow` again when flai serve is not running. The running flai serve's status already holds `ImportProjects` and `Unserved`. `Host.Settings` (`flai/internal/hostapi/writes.go`) is `func(root string) any`, so `settings.get` has no context to record phases with.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/cmd/serve_actions.go` | declared; `hostSettings` and `servedView` |
| `flai/cmd/serve_project.go` | declared; `servedBelow`, `importRootOf`, `listServedProjects` |
| `flai/internal/serve/below.go` | declared; `FindBelow`, where a helper that answers a root's import folder from one found list may go |
| `flai/cmd/serve_project_test.go` | co-change (9%); layout: it already tests `servedView`, and criterion 1's test of three served projects under two import folders goes there |
| `flai/cmd/serve_actions_test.go` | co-change (65%); layout: its tests call `hostSettings(root)`, whose signature takes a context for criterion 2 |
| `flai/internal/hostapi/hostapi.go` | co-change (15%); layout: the `settings.get` handler, where its phases start |
| `flai/internal/hostapi/writes.go` | co-change (37%); layout: `Host.Settings`, which takes a context |
| `flai/internal/hostapi/writes_test.go` | co-change (32%); layout: `TestSettingsGetSaysWhatMayChange` builds a `Host.Settings` |
| `design/system/server-performance.md` | design: criterion 3 records the before and after under S-0596's section |

Not taken from `flai touches suggest`: `docs/users/flai-reference.md`, `docs/users/flai.md`, `docs/operators/settings.md`, and `docs/operators/index.md`, since no command, flag, or setting changes; `flai/internal/serve/agents.go`, `flai/internal/serve/serve.go`, `flai/cmd/serve.go`, `flai/internal/config/config.go`, and `flaiover/src/lib/server/agent.ts`, since the fix reads the serve status and the configuration as they are; `flai/internal/hostapi/settings.go`, which holds the setting writes, not `settings.get`; `design/system/flai-cli.md`, whose account of `settings.get`'s `host.projects` stays true. No folder touch.

`after`: S-0596, now done. It wrote the section of `design/system/server-performance.md` that the goal links, cause 15 and its figures, which criterion 3 measures against and records under.

S-0603 (cause 12) also changes `flai/cmd/serve_actions.go`, and S-0604 and S-0606 change `flai/internal/hostapi/hostapi.go`. None writes output this story reads, so none is in `after`; the trial merge judges the overlap.

Figures:

- Forecast 30m, against flai's 8m (79 s per unit of size over 20 done improvement stories on claude-opus-5-5 in the small band, times size 6). The size counts three touches, but the change threads a context through the host API and its tests and ends in a measurement on the operator's host. S-0606 forecasts 45m for two mechanisms and a measurement; this one is about two thirds of that.
- Delivery 2026-11-08T00:33:00Z: flai's delivery for its queue position, 9th behind S-0591 to S-0604 with an in-progress limit of 5, moved by the 22 minutes the duration adds.
- Cost of delay: no value yet. Neither S-0605 nor E-0033 has inputs, so `flai cod S-0605` refuses. The plan's thread on S-0605 asks for `penalty_per_week: 20`, the figure on TH-0888; the value follows once an input is set.

### Accepted by the orchestrator

- Verified: 8793f269a5f7f4d9a7d1fc0bd7467493fa43f5a2
- At: 2026-11-08T01:13:20Z

Verdict: pass at 8793f269a5f7f4d9a7d1fc0bd7467493fa43f5a2
- 1: flai/cmd/serve_actions.go, flai/cmd/serve_project.go, flai/internal/serve/below.go, flai/cmd/serve_project_test.go, flai/internal/serve/below_test.go
- 2: flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi.go, flai/cmd/serve_actions.go, flai/cmd/serve_settings.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes_test.go
- 3: design/system/server-performance.md
