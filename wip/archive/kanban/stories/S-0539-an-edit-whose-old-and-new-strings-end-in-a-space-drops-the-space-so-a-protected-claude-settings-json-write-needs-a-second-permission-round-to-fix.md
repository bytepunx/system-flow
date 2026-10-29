---
id: S-0539
type: story
nature: improvement
title: An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix
status: done
owner: alex
created: 2026-10-29T02:01:44Z
updated: 2026-10-29T05:16:30Z
transitions:
  - to: ready
    at: 2026-10-29T03:56:58Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-29T03:57:27Z
    by: agent-S-0539
  - to: review
    at: 2026-10-29T04:58:57Z
    by: agent-S-0539
  - to: done
    at: 2026-10-29T05:16:30Z
    by: orchestrator
tags: [flai, mcp, permission-prompt]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, docs/users/flai.md, design/system/agent-adapters.md, design/issues/I-0223-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md, design/issues/summary.md, design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3700
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 2
      hand_edits: 1
      work: 40
  models:
    - model: claude-opus-5-5
      input: 96
      output: 17624
      cache_read: 7327653
      cache_write: 222057
      cost: 3.5949
  strategic:
    - kind: orchestrator
      seconds: 2797
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 159
          output: 2636
          cache_read: 28238695
          cache_write: 47862
          cost: 6.238
        - model: claude-sonnet-5-5
          input: 12
          output: 76
          cache_read: 140024
          cache_write: 47412
          cost: 0.1292
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-29T02:01:44Z
  value: 37.5
  by: planner-S-0539
  at: 2026-10-29T03:02:42Z
forecast:
  duration: 1h
  delivery: 2026-10-29T20:27:00Z
  basis: "Its own forecast of 1h; 35th in the pull order with an in-progress limit of 5, behind S-0417, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0508, S-0509, S-0510, S-0511, S-0512, S-0515, S-0516, S-0517, S-0518, S-0520, S-0523, S-0524, S-0527, S-0529, S-0418, S-0419, S-0420, S-0421, S-0422 and S-0423."
  by: flai
  at: 2026-10-29T03:57:15Z
finalized:
  by: orchestrator
  at: 2026-10-29T03:56:46Z
---
# S-0539 An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix

## Goal

This story remediates [I-0223](../../../design/issues/I-0223-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md), "An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0223 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0223 is closed with `flai issue close I-0223 --reason` saying what fixed it

## Tasks
- T-2320 permission_prompt refuses at once an edit whose old string ends in a space or tab and whose new string does not
- T-2321 Document permission_prompt's refusal of an edit whose new string lost the old string's trailing space
- T-2322 Close I-0223 naming permission_prompt's refusal as its fix

## Notes

Cost of delay inputs set by flai from I-0223. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T22:40:11Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the instance. TH-0712, the first of S-0416's permission threads, shows what `permission_prompt` received: an old string ending in `exit 0;` and a space, and a new string ending in `exit 0;` with no space. So the space was gone before flai saw the call: Claude Code drops trailing whitespace from an Edit's `new_string` and keeps it on `old_string`. flai passes `input` back as `updatedInput` unchanged (`allow` in `flai/internal/mcpserver/permission.go`), so flai cannot restore the space, and the operator approved a change that looked right. The remedy is to refuse such an edit at once in `permission_prompt`, before any thread and under auto-approve too: an `Edit`, or an edit of a `MultiEdit`, whose `old_string` ends in a space or a tab while its `new_string` does not. The refusal tells the agent to end both strings on a non-space character. That costs one immediate retry instead of a wrong write and a second permission round. The first task builds it and reproduces S-0416's instance in a test.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/mcpserver/permission.go` | layout: `permission_prompt`'s decision and its tool description |
| `flai/internal/mcpserver/permission_test.go` | layout: the test that reproduces the instance |
| `docs/users/flai.md` | design: the documentation convention; its permission section and MCP tool table describe what `permission_prompt` refuses at once |
| `design/system/agent-adapters.md` | design: its section on `permission_prompt` lists what it refuses |
| `design/issues/I-0223-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai touches suggest S-0539` with the two permission files listed `flai/internal/mcpserver/folder.go` (2 of 7 commits, 29%). It is left out: the tool's description is a constant in `permission.go`, and `folder.go` only registers it. No folder touch is kept.

Holds: S-0532, in progress, claims `flai/internal/mcpserver` and `docs`, so this story is held behind it until it moves to review. S-0519, in review, touches `permission.go` too and holds nothing. Neither is named in `after`: this story needs neither one's change on main.

Figures:

- Forecast: `flai forecast` gave 3m from size 2 (2 criteria, 0 touches) at 80 s a unit over 13 small improvement stories. That figure was taken before any touches were written. It is raised to 1h for one Go check with its tests, two documents, the issue's close, and the close-out's test tiers. The delivery stays at flai's 2026-10-29T22:24:00Z, which the 36 stories ahead in the pull order set more than this story's own time does.
- Cost of delay: `flai cod` gives 37.50 USD a week, from the input flai set from I-0223: 15m lost per 168h cycle at 150 USD an hour. It stands, since the issue has one instance.

### Accepted by the orchestrator

- Verified: 58c1387f8080596f44cb8d2724f3b9aef3e2cd89
- Batch: 20261029-045945
- Enqueued verified: 58c1387f8080596f44cb8d2724f3b9aef3e2cd89
- At: 2026-10-29T05:16:30Z

Verdict: S-0539 meets both criteria at 58c1387f; flai verify passed with integration and smoke left to the queue, and every changed file is under the story's touches.
- 1: flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, docs/users/flai.md, design/system/agent-adapters.md
- 2: design/issues/I-0223-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md, design/issues/summary.md
