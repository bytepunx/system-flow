---
id: S-0616
type: story
nature: remediation
title: TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input
status: done
owner: alex
created: 2026-11-07T23:39:18Z
updated: 2026-11-09T22:22:33Z
transitions:
  - to: ready
    at: 2026-11-09T10:22:28Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T10:22:38Z
    by: agent-S-0616
  - to: review
    at: 2026-11-09T10:43:49Z
    by: agent-S-0616
  - to: in-progress
    at: 2026-11-09T16:57:24Z
    by: alex
  - to: review
    at: 2026-11-09T21:52:58Z
    by: alex
  - to: in-progress
    at: 2026-11-09T21:53:09Z
    by: flai
  - to: review
    at: 2026-11-09T22:09:30Z
    by: agent-S-0616
  - to: done
    at: 2026-11-09T22:22:33Z
    by: orchestrator
tags: [flai, template]
touches: [flai/cmd/new_test.go, template/root/design/workflows/formulas/README.md, design/workflows/formulas/README.md, design/system/template.md, template/CHANGELOG.md, design/issues/I-0249-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md, design/issues/summary.md, design/issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2266
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 7
      hand_edits: 2
      work: 300
  models:
    - model: claude-opus-5-5
      input: 630
      output: 41500
      cache_read: 98610365
      cache_write: 644764
      cost: 25.7127
  strategic:
    - kind: planner
      seconds: 222
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 17689
          cache_read: 2978465
          cache_write: 149480
          cost: 2.1455
    - kind: orchestrator
      seconds: 89
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 168
          cache_read: 801407
          cache_write: 126162
          cost: 0.2431
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T23:39:18Z
  value: 25
  by: planner-S-0616
  at: 2026-11-09T10:20:35Z
forecast:
  duration: 1h
  delivery: 2026-11-09T11:33:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 5, behind S-0608, S-0610 and S-0615."
  by: flai
  at: 2026-11-09T10:22:32Z
finalized:
  by: orchestrator
  at: 2026-11-09T10:22:22Z
---
# S-0616 TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input

## Goal

This story remediates [I-0249](../../../design/issues/I-0249-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md), "TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0249 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0249 is closed with `flai issue close I-0249 --reason` saying what fixed it

## Tasks
- T-2720 TestRenderPrototypeTemplate looks for template syntax only in files the renderer rendered or changed
- T-2721 The formulas README shows an input placeholder in braces again, and template.md says copied files may hold them
- T-2722 Close I-0249 with the fix named

## Notes

Cost of delay inputs set by flai from I-0249. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T02:50:10Z, 5.9 days before this story; under one cycle counts as one).

### Planning

Proposed solution. `flai/internal/template/render.go` renders only the files whose name ends in the template's `render.suffix` (`.tmpl`), with Go `text/template`, and copies every other file byte for byte. So a `{{` in a copied file is content, not syntax left unrendered. But `TestRenderPrototypeTemplate` in `flai/cmd/new_test.go` refuses `{{` in every file of the rendered project. S-0424 and S-0554 worked around it for two cases: `withoutFormulaInputs` drops a formula's declared input placeholders, and `withoutRoleStart` drops a role file's `start` prompt. Any other copied file still fails, such as the formulas README that S-0421 had to reword. The fix is to check for `{{` only where rendering could leave it: in each file rendered from a `.tmpl` source, and in each file that has no source in `template/root` or does not match its source byte for byte. A file that matches its verbatim source exactly is passed over. The two workaround helpers and their tests then go. The renderer does not change.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/new_test.go` | layout | The test and both helpers live here (T-2720). |
| `template/root/design/workflows/formulas/README.md` | co-change, design | The instance in I-0249. It shows `{{name}}` again to prove the fix on the real template (T-2721). |
| `design/workflows/formulas/README.md` | co-change (29%) | This project's copy of that README, identical today, kept in step (T-2721). |
| `design/system/template.md` | co-change (18%), design | Its `render.suffix` line says which files are rendered. It gains that the others are copied byte for byte and may hold `{{` (T-2721). |
| `template/CHANGELOG.md` | co-change (24%) | A patch line for the README change under `## Unreleased` (T-2721). |
| `design/issues/I-0249-…md` | design | Closed with `flai issue close` (T-2722). |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it (T-2722). |

`flai touches suggest` also listed `flai/cmd/new.go`, `docs/users/flai.md`, and `design/system/flai-cli.md`. They are left out because no command or renderer behaviour changes. There are no folder touches and no `after`.

Figures:

- Forecast: `flai forecast` gave 5m, from the median 144 s per size unit over 17 small remediation stories, at size 2 (two criteria, no touches). I raised it to 1h. The story has three tasks over `flai/` and `template/`. Its close-out runs the integration and smoke tiers, and `TestRenderPrototypeTemplate` runs only in integration (it skips under `-short`). Delivery is set to 12:30Z because the story waits behind the stories in progress.
- Cost of delay: 25 USD a week, as `flai cod` worked it out from the inputs flai set from I-0249: 10m lost per 168h cycle at 150 USD an hour. It stands. The issue's one instance was worked around in S-0421, and S-0424 and S-0554 narrowed the cause. The remaining cost is the next workaround someone has to find.
- 2026-11-09T16:57:24Z: moved to in-progress: Merge failed, need to rebase
- 2026-11-09T21:53:09Z: moved to in-progress: the rebase of story/S-0616 onto main stopped on conflicts in design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, which the acceptance of S-0616 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0616

### Accepted by the orchestrator

- Verified: 377709c9aeb87d31dd254995152c7f2971ea0ca1
- At: 2026-11-09T22:22:33Z

Verdict: pass; both criteria met at 377709c, close-out passed every step there.
- 1: flai/cmd/new_test.go, template/root/design/workflows/formulas/README.md, design/workflows/formulas/README.md, design/system/template.md, template/CHANGELOG.md
- 2: design/issues/I-0249-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md, design/issues/summary.md
