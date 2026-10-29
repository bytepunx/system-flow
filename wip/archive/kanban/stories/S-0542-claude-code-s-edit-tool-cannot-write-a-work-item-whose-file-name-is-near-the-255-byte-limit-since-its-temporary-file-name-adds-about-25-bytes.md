---
id: S-0542
type: story
nature: improvement
title: Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes
status: done
owner: alex
created: 2026-10-29T15:14:37Z
updated: 2026-10-29T15:49:53Z
transitions:
  - to: ready
    at: 2026-10-29T15:25:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-29T15:25:11Z
    by: agent-S-0542
  - to: review
    at: 2026-10-29T15:33:39Z
    by: agent-S-0542
  - to: done
    at: 2026-10-29T15:49:53Z
    by: orchestrator
tags: []
touches: [flai/internal/template/render.go, flai/internal/template/template_test.go, flai/internal/workitem/repo.go, flai/internal/workitem/create_test.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/internal/adr/adr.go, flai/internal/digest/write.go, design/system/work-hierarchy.md, design/system/documentation-standard.md, docs/users/flai.md, design/issues/I-0231-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 522
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 2
      hand_edits: 1
      work: 24
  models:
    - model: claude-opus-5-5
      input: 148
      output: 35936
      cache_read: 5353458
      cache_write: 294414
      cost: 3.7131
  strategic:
    - kind: planner
      seconds: 336
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 85
          output: 1342
          cache_read: 4834730
          cache_write: 160061
          cost: 1.2476
    - kind: orchestrator
      seconds: 893
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 154
          output: 2486
          cache_read: 56910935
          cache_write: 103894
          cost: 12.5727
        - model: claude-sonnet-5-5
          input: 12
          output: 63
          cache_read: 158573
          cache_write: 46629
          cost: 0.1414
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-29T15:14:37Z
  value: 7.5
  by: planner-S-0542
  at: 2026-10-29T15:23:08Z
forecast:
  duration: 45m
  delivery: 2026-10-30T07:13:00Z
  basis: "flai forecast's 29m (89 s per unit over 103 done improvement stories, size 19) raised to 45m for four tasks each closed with a commit, sync, and check, and the Go test tiers of six packages in the close-out; delivery is flai's 31st-in-pull-order date plus the 16m."
  by: planner-S-0542
  at: 2026-10-29T15:23:08Z
finalized:
  by: orchestrator
  at: 2026-10-29T15:24:56Z
---
# S-0542 Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes

## Goal

This story remediates [I-0231](../../../design/issues/I-0231-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md), "Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0231 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0231 is closed with `flai issue close I-0231 --reason` saying what fixed it

## Tasks
- T-2333 A work item's file name cuts its title's slug at 80 bytes on a word, through one helper in the template package
- T-2334 Issue, thread, and conversation file names cut their slug at 80 bytes with the helper, and the ADR and digest cuts move onto it
- T-2335 The design and the user docs say a file name flai makes from a title cuts its slug at 80 bytes on a word
- T-2336 Close I-0231 with the reason that flai now cuts a title's slug at 80 bytes in every file name it makes

## Notes

Cost of delay inputs set by flai from I-0231. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-29T14:15:00Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0231's one instance: every file name flai makes from a title cuts the title's slug at 80 bytes on a word, as `flai adr new` (`flai/internal/adr/adr.go`) and `flai digest new` (`flai/internal/digest/write.go`) already do. One helper in `flai/internal/template/render.go` does the cut, and the work item name (`workitem.FileName`, which `flai story new`, `task new`, `epic new`, a retitle in `itemedit`, and an experiment's expected results document use), the issue name, the thread name, and the conversation name call it; the ADR and digest cuts move onto it. A name is then at most the ID, a dash, 80 bytes, and `.md`, about 92 bytes, far under the 255-byte limit with Claude Code's `<name>.tmp.<pid>.<hex>` added. Today 8 open files under `wip/kanban`, `design/issues`, and `wip/threads` are over 225 bytes, the longest 239. `flai check`'s `item.filename`, `issues.filename`, and `threads.filename` only ask that a name start with its ID, so a cut name passes them, and experiment documents are found by their ID prefix (`experiment.Find`), so an existing results document is still found.

Touches, all files:

| Touch | From |
|-------|------|
| `flai/internal/template/render.go`, `template_test.go` | layout: `Slug` lives there; the cut and its test |
| `flai/internal/workitem/repo.go`, `create_test.go` | goal and layout: `FileName` names every work item; the test reproduces I-0231 |
| `flai/internal/issues/issues.go`, `issues_test.go` | co-change (24%) and layout: the issue name; I-0231's own name is 239 bytes |
| `flai/internal/threads/threads.go`, `threads_test.go` | co-change (21%) and layout: the thread name; TH-0465's is 238 bytes |
| `flai/internal/messages/messages.go`, `messages_test.go` | layout: the conversation name, built as the thread's is |
| `flai/internal/adr/adr.go`, `flai/internal/digest/write.go` | layout: the two cuts that exist, moved onto the helper |
| `design/system/work-hierarchy.md` | co-change (12%) and design: its `File name:` line |
| `design/system/documentation-standard.md` | design: § Naming says how items are named |
| `docs/users/flai.md` | co-change (26%): user docs change with a command's behaviour (documentation.md) |
| `design/issues/I-0231-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

Left out of what `flai touches suggest` listed: `design/system/flai-cli.md` (35%), `flai/cmd/root.go`, and `flai/internal/check/check.go`, since no command's flags or check rule change; `flaiover/src/lib/server/repo.ts`, since the dashboard reads item paths from flai and builds none from a title. No folder touch.

Forecast: `flai forecast` gave 3m from size 2 (2 criteria, 0 touches), before the touches were written. Adjusted to 1h30m: four tasks across six Go packages, each with a test, three documents, and a close-out that runs the Go test tiers for those packages. Delivery is `flai forecast`'s, the story being 31st in the pull order with the in-progress limit of 5.

Cost of delay: `flai cod` gives 7.50 USD a week from the input flai set from I-0231 (3m lost per 168h cycle at 150 USD an hour). It stands: the input counts one occurrence, and each new long title is another, so the value may rise as I-0231 is bumped.

Assumptions: the cut applies to names made from now on; the 8 open over-long files keep their names until retitled (a retitle names the file with the cut), since renaming them would break the links to them. No ADR: the rule extends to every title-named file what ADR-0191 and `flai adr new` already do.

### Accepted by the orchestrator

- Verified: a2c7a0044225715eca8ea360c59ac4a0931a309a
- Batch: 20261029-153625
- Enqueued verified: a2c7a0044225715eca8ea360c59ac4a0931a309a
- At: 2026-10-29T15:49:53Z

Verdict: S-0542 meets both criteria at a2c7a004; flai verify passed with integration and smoke left to the queue, all 17 changed files are under the story's touches, and the ADR and digest cuts are unchanged.
- 1: flai/internal/template/render.go, flai/internal/template/template_test.go, flai/internal/workitem/repo.go, flai/internal/workitem/create_test.go, flai/internal/issues/issues.go, flai/internal/threads/threads.go, flai/internal/messages/messages.go, flai/internal/adr/adr.go, flai/internal/digest/write.go, design/system/work-hierarchy.md, design/system/documentation-standard.md, docs/users/flai.md
- 2: design/issues/I-0231-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md, design/issues/summary.md
