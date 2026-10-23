---
id: S-0468
type: story
nature: improvement
title: Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five
status: review
owner: alex
created: 2026-10-21T11:58:08Z
updated: 2026-10-23T05:13:26Z
transitions:
  - to: ready
    at: 2026-10-22T20:51:26Z
    by: alex
  - to: in-progress
    at: 2026-10-23T03:33:07Z
    by: agent-S-0468
  - to: review
    at: 2026-10-23T05:13:26Z
    by: agent-S-0468
tags: [cli, template]
topics: [planning]
touches: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/planning_test.go, flai/internal/itemedit/folder_test.go, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, ".claude/agents/planner.md", template/root/.claude/agents/planner.md, template/template.yaml, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md, design/issues/summary.md, flai/cmd/edit.go, flai/cmd/check.go, flai/cmd/check_stats_test.go, design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6007
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 1
      hand_edits: 4
      work: 82
  models:
    - model: claude-opus-5-5
      input: 538
      output: 182984
      cache_read: 35802047
      cache_write: 1167610
      cost: 17.6108
  strategic:
    - kind: orchestrator
      seconds: 18072
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 170
          output: 2578
          cache_read: 46419500
          cache_write: 71194
          cost: 10.2655
        - model: claude-sonnet-5-5
          input: 12
          output: 49
          cache_read: 213243
          cache_write: 75897
          cost: 0.2001
    - kind: analyzer
      seconds: 113
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 97
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
cost_of_delay:
  inputs:
    time_lost_per_cycle: 44h
    by: flai
    at: 2026-10-21T11:58:08Z
  value: 6600
  by: planner-S-0468
  at: 2026-10-22T21:04:06Z
forecast:
  duration: 2h
  delivery: 2026-10-23T00:57:00Z
  basis: "Its own forecast of 2h; 1st in the pull order with an in-progress limit of 5, behind S-0466."
  by: flai
  at: 2026-10-22T22:37:19Z
finalized:
  by: alex
  at: 2026-10-22T20:51:09Z
---
# S-0468 Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five

## Goal

This story remediates [I-0152](../../../design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md), "Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0152 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0152 is closed with `flai issue close I-0152 --reason` saying what fixed it

## Tasks
- T-2055 An ADR records the remedy for I-0152: a task takes no folder touch, a planner's folder touch on a story needs its Planning note, and a kept folder or an after hub is a check finding
- T-2056 flai refuses a task a folder touch through every write path, and workitem tells a folder of the tree from a file, leaving out the shared paths
- T-2057 The planner's folder touch on a story is refused unless the story's Planning note names it, through flai edit, flai touches, item_new, and item_edit
- T-2058 flai check reports touches.folder on a story in progress whose claim keeps a folder and wip.after-hub on a story three or more open stories wait for, and a scoped check drops another story's
- T-2059 The conventions, the planner's prompt and definition, the user docs, and the template say a task names files, a planner's kept folder needs its note, after names a needed change, and the two findings
- T-2060 I-0152 is closed with what fixed it and the story's criteria are verified and ticked

## Notes

The analyzer found I-0152: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0152. time_lost_per_cycle 44h carried over from I-0152's Impact section. Its Impact time_lost_per_cycle was taken rather than the 44h derived from its cost and count.

Strategic usage carried over from I-0152: analyzer 1m53s, 1.2510 USD.

### Planning

Planned by planner-S-0468 on 2026-10-22. The remedy the tasks build, proposed from I-0152's instances in the plan's thread: a task takes no folder touch, so that ADR-0096's narrowing always has files to narrow to; the planner's folder touch on a story is refused unless the story's `### Planning` names it; `flai check` reports a story in progress whose claim keeps a folder (`touches.folder`) and a story three or more open stories wait for (`wip.after-hub`); and the planner's convention and prompt say `after` names a story whose change the story needs on main, never an order of preference.

Touches, by source. Declared: none; the story declared no touches. Layout, the code the remedy changes: `flai/internal/workitem/create.go` and `create_test.go` (`CleanTouches` and `CheckTouch`, where every touches write passes), `hold.go` and `hold_test.go` (the claim, and the folder test the refusals and the check share), `flai/internal/itemedit/itemedit.go` and `planning_test.go` (`Apply`, behind `flai edit` and `item_edit`), `flai/internal/itemnew/itemnew.go` and `itemnew_test.go` (the planner's checked creation), `flai/cmd/touches.go` and `touches_test.go` (`flai touches` saves without `Apply`), `flai/internal/mcpserver/items_write.go` and `items_write_test.go`, `flai/internal/check/check.go`, `check_test.go`, `scope.go`, and `scope_test.go` (the two rules and the close-out's scope), `flai/internal/harness/harness.go` and `harness_test.go` (`planPrompt`). Design: `design/adrs/README.md`, `design/system/strategic-agents.md` § Enriching a story, `design/system/workflow.md` § Branches and collisions, `design/system/flai-cli.md` (the check rules), the conventions `strategic-agents.md` and `work-management.md` with their template copies, `.claude/agents/planner.md` and its template copy (the design says they say the same as `planPrompt`; they lack the folder rule), `template/template.yaml` and `template/CHANGELOG.md` (a template release), `design/issues/I-0152-…md` and `summary.md` (the close). Co-change, from `flai touches suggest` over seven seed files: `harness_test.go` 38%, `flai-cli.md` 21%, `docs/users/flai.md` 21%, `template/CHANGELOG.md` 19%, `check_test.go` 15%, the template's `strategic-agents.md` 13%, `template.yaml` 12%, `work-management.md` 10%, all kept; `docs/users/flai-reference.md` 8% kept in case `flai touches`'s help text changes. Left out: `serve/agents.go`, `cmd/items.go`, `cmd/edit.go`, `cmd/root.go`, `mcpserver/server.go`, `manifest.go`, `hostapi/writes.go`, and `docs/operators` (no setting, flag, or serve behaviour changes; `items.go` and `edit.go` reach the rule through `Repo.Create` and `Apply`).

Folder touch kept: `design/adrs`, because the first task's ADR gets its file name from `flai adr new`, which no task can name yet; it lies inside `claims.shared`, so it holds no ready story.

Forecast: flai gave 55m, 89 s per unit over 87 done improvement stories on claude-opus-5-5 in the large band, times size 37. Raised to 2h: the size counts 35 touches, half of them test and document files, while the work is three refusals and two check rules across six Go packages with tests, an ADR, conventions, and a template release, in six tasks over five layers. Delivery 2026-10-22T23:30:00Z: now plus 2h times the cycle factor 1.15, with a lane free (three of five in progress).

Cost of delay: 6600 USD a week, as flai works it out from the input on the story, 44h lost per 168h cycle at 150 USD an hour, one cycle a week. It stands: no revenue or penalty input is set, and the input is the operator's.
