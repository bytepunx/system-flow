---
id: S-0544
type: story
nature: improvement
title: Make the dashboard and server versions easy to see
status: done
owner: alex
created: 2026-10-29T17:39:08Z
updated: 2026-10-29T18:02:07Z
transitions:
  - to: ready
    at: 2026-10-29T17:39:15Z
    by: alex
  - to: in-progress
    at: 2026-10-29T17:39:23Z
    by: agent-S-0544
  - to: review
    at: 2026-10-29T17:50:56Z
    by: agent-S-0544
  - to: done
    at: 2026-10-29T18:02:07Z
    by: alex
tags: []
touches: [flaiover/src, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, flaiover/src/routes/api/dashboard/+server.ts, flaiover/src/routes/api/dashboard/dashboard.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 701
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 2
      hand_edits: 1
      work: 47
  models:
    - model: claude-opus-5-5
      input: 114
      output: 23827
      cache_read: 8310541
      cache_write: 191506
      cost: 3.6712
  strategic:
    - kind: orchestrator
      seconds: 133
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 54
          output: 376
          cache_read: 1364447
          cache_write: 173098
          cost: 0.3408
        - model: claude-sonnet-5-5
          input: 12
          output: 105
          cache_read: 140046
          cache_write: 37658
          cost: 0.1225
forecast:
  duration: 1h
  delivery: 2026-10-29T18:45:00Z
  basis: flai forecast's 4m (80 s per unit over 13 small improvement stories, size 3) raised to 1h for two tasks that each change a Svelte component, an API route, their tests, and two docs, and for the close-out's flaiover tiers
  by: planner-S-0544
  at: 2026-10-29T17:42:53Z
---
# S-0544 Make the dashboard and server versions easy to see

## Goal

Right now, the host page doesn't show the dashboard version and you tend to only see that new versions are available after a push. It should be very easy to see what versions are running.

## Acceptance criteria
- [x] Put a label on the nav rail under the license that shows the server and dashboard versions (newline separated)
- [x] Show the versions of the current dashboard and the available dashboard version when the user checks, not just that a new version is available

## Tasks
- T-2341 The nav rail shows the flai and dashboard versions under License
- T-2342 The Dashboard area shows the running dashboard version, and Check for updates names the running and the available version

## Notes

### Planning

Touches:

- `flaiover/src`: declared by the operator, kept. It is a folder touch: the story's claim replaces it with its tasks' file touches (ADR-0096), and the agent may still add a file there, such as a client module declaring `__FLAIOVER_VERSION__`, that no task names yet.
- `NavRail.svelte` and its test: layout. The License label is the rail's foot (`flaiover/src/lib/components/NavRail.svelte`), so criterion 1 lands there. Co-change also lists `design/components/flaiover/src/lib/components/NavRail.examples.json`; it is left out, since the label reads `hostFlai` state, not props.
- `hostflai.svelte.ts`, `routes/api/agent/+server.ts` and its test: layout. The rail learns the connected flai's version from `GET /api/agent`, and T-2341 adds the dashboard's own version there.
- `HostPanel.svelte` and its test, `routes/api/dashboard/+server.ts` and its test: layout. The Check for updates result is rendered in `HostPanel.svelte` and today says only whether an update is available.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`: co-change (16% and 20% of the commits that changed `flaiover/src`) and design. They describe the rail's foot and the host page's Dashboard area.
- No `flai/` path: `dashboard.check` already answers `release` and `image` for the available version, and the running version comes from the dashboard itself. A change on the flai side would show only once flai is released and installed.

Figures:

- Forecast 1h, not flai's 4m: the history's median rests on tiny stories, and this one has two tasks, each with a component, an API route, tests, and two docs, plus the close-out's flaiover tiers.
- Cost of delay: no inputs and no epic, so `flai cod` refuses it. Asked on TH-0756, recommending `time_lost_per_cycle: 10m` (25 USD/week).
