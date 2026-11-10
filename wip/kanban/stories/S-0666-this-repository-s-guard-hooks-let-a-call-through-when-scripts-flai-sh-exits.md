---
id: S-0666
type: story
nature: improvement
title: This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would
status: in-progress
owner: alex
created: 2026-11-10T08:34:38Z
updated: 2026-11-10T11:30:23Z
transitions:
  - to: ready
    at: 2026-11-10T10:57:24Z
    by: alex
  - to: in-progress
    at: 2026-11-10T11:14:16Z
    by: agent-S-0666
blocked:
  - from: 2026-11-10T11:29:23Z
    until: 2026-11-10T11:30:01Z
    reason: awaiting the operator on TH-1051
  - from: 2026-11-10T11:30:23Z
    reason: awaiting the operator on TH-1051
tags: [guard, scripts, harness]
topics: [harness, cli]
touches: [design/adrs, design/adrs/README.md, scripts/flai.sh, flai/cmd/guard_test.go, ".opencode/plugins/flai-guard.js", flai/tests/integration/opencode_test.go, scripts/README.md, docs/users/flai.md, design/system/agent-adapters.md, docs/operators/index.md, design/system/strategic-agents.md, design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md, design/issues/summary.md, design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 926
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 2
      work: 44
  models:
    - model: claude-opus-5-5
      input: 210
      output: 57341
      cache_read: 14503681
      cache_write: 478396
      cost: 7.398
  strategic:
    - kind: planner
      seconds: 292
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 318
          output: 8609
          cache_read: 27326488
          cache_write: 536278
          cost: 1.8965
        - model: claude-opus-5-5
          input: 871
          output: 86111
          cache_read: 68935829
          cache_write: 1044206
          cost: 23.6176
    - kind: orchestrator
      seconds: 66
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 30
          output: 209
          cache_read: 1255404
          cache_write: 131023
          cost: 0.3799
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-11-10T10:55:48Z
forecast:
  duration: 45m
  delivery: 2026-11-10T12:39:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 3, behind S-0665, S-0667 and S-0708."
  by: flai
  at: 2026-11-10T11:05:34Z
finalized:
  by: alex
  at: 2026-11-10T10:57:18Z
---
# S-0666 This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would

## Goal

This story remediates [I-0275](../../../design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md), "This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0275 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0275 is closed with `flai issue close I-0275 --reason` saying what fixed it

## Tasks
- T-3061 Write the ADR refining ADR-0235: scripts/flai.sh guard runs an older bin/flai or the installed flai when the build fails, and exits 127 with neither
- T-3066 scripts/flai.sh guard falls back to an older bin/flai or the installed flai when go build fails, and exits 127 with neither, with a test that makes the build fail
- T-3069 This repository's OpenCode flai-guard plugin refuses a call when scripts/flai.sh exits 126 or 127 in a session flai started
- T-3070 Document scripts/flai.sh's guard fallback on a failed build and the OpenCode plugin's 126 and 127 refusal
- T-3071 Close I-0275 with what fixed it

## Notes

### Planning

Proposed fix, from I-0275's one instance (S-0567's live planner run) and ADR-0235. `scripts/flai.sh` runs `scripts/flai-build.sh` when `go` is present and `bin/flai` is stale; a failed `go build` exits 1 under `set -eu`, and the hooks in `.claude/settings.json` allow every exit but 2, 126, and 127. The fix keeps the hooks as they are and makes `scripts/flai.sh guard` always reach a verdict: when the build fails, it runs the existing `bin/flai`, older than the sources, or else the installed `flai`, each with a warning on standard error, and exits 127 when neither is there, so the hook refuses. Other commands keep failing on a failed build, so an agent still sees a broken build. ADR-0235 says a failed build fails open and that `scripts/flai.sh` never runs an older installed flai where `go` is present, so a new ADR refining ADR-0235 records the change first. This repository's `.opencode/plugins/flai-guard.js` runs `scripts/flai.sh` too and treats only a spawn error as a missing guard, so it passes an exit 127; it gets the same 126 and 127 rule as the Claude Code hooks.

Touches, and where each came from (`flai touches suggest` listed nothing: the story declared no touches):

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs` | design | The new ADR refining ADR-0235. Its number and slug are not known until `flai adr new` writes it, and other open stories write ADRs meanwhile, so no task can name the file yet; the folder stays on the story for that one file. |
| `design/adrs/README.md` | layout | `flai adr new` adds the ADR's row to the generated index. |
| `scripts/flai.sh` | layout | The fallback on a failed build under `guard`. |
| `flai/cmd/guard_test.go` | layout | `TestTheGuardHooksRefuseACallWhenTheGuardCannotStart` runs both hook commands with no `go` on `PATH`; a fake `go` that exits 1 reproduces I-0275. |
| `.opencode/plugins/flai-guard.js` | layout | This repository's OpenCode guard, which runs `scripts/flai.sh` and passes an exit 127. |
| `flai/tests/integration/opencode_test.go` | layout | `TestFlaiGuardPluginUnderNode` loads only the template's plugin; it loads this repository's too. |
| `scripts/README.md` | design | The `flai.sh` row describes its fallback. |
| `docs/users/flai.md` | design | `## Sub-agents` describes the 127 refusal and the installed-flai fallback. |
| `design/system/agent-adapters.md` | design | `### The guard` holds the exit table; `### A container harness of flai's own` describes the fallback. |
| `docs/operators/index.md` | design | `### Strategic agents in a container (S-0509)` says `scripts/flai.sh` needs `go`. |
| `design/system/strategic-agents.md` | design | `### The live run (S-0567)` records I-0275 as open. |
| `design/issues/I-0275-...md`, `design/issues/summary.md` | declared by the goal | `flai issue close` writes both. |

No `after`: the story reads nothing another open story writes.

Forecast: `flai forecast` gave 26m, from 103 s per unit of size over 155 done large improvement stories, times size 15. Raised to 45m: the five tasks run in three layers one after another, the ADR first, and the `flai/` touches make the close-out run the integration and smoke tiers. The delivery moved 19m later to match, to 2026-11-10T14:15:00Z.

Cost of delay: no value yet. The story has no inputs and no epic, so `flai cod` cannot work one out. TH-1039 asks the operator for an input, recommending `time_lost_per_cycle: 5m`.

Tasks, in three layers: T-3061 (the ADR); then T-3066 (`scripts/flai.sh` and its test) and T-3069 (the OpenCode plugin and its test), which share no path; then T-3070 (docs) and T-3071 (close I-0275), which share no path.
