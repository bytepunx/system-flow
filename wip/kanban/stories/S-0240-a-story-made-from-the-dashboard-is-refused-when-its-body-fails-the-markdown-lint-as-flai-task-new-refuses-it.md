---
id: S-0240
type: story
nature: remediation
title: A story made from the dashboard is refused when its body fails the markdown lint, as flai task new refuses it
status: in-progress
owner: arobson
created: 2026-10-02T12:42:51Z
updated: 2026-10-02T16:50:12Z
transitions:
  - to: ready
    at: 2026-10-02T16:24:25Z
    by: alex
  - to: in-progress
    at: 2026-10-02T16:48:17Z
    by: agent-S-0240
tags: [dashboard, cli]
topics: [dashboard, markdown]
touches: [flai/internal/mdlint, flai/internal/itemnew, flai/internal/hostapi, flaiover/src/lib/server, design/issues, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 580
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 164
      output: 1098
      cache_read: 5614426
      cache_write: 214563
      cost: 2.3504
---
# S-0240 A story made from the dashboard is refused when its body fails the markdown lint, as flai task new refuses it

## Goal

[I-0055](../../../design/issues/I-0055-a-story-made-from-the-dashboard-reached-main-with-a-markdown-lint-error-flai-task-new-would-have-refused.md): S-0231, made from the dashboard on 2026-10-02, reached `main` with its criteria list indented one space (MD007), which stopped S-0193's close-out at the markdown lint. `flai task new` and `flai story new` refuse such a body with the rule and the line (ADR-0061), so the dashboard's path for a new item, or the edit that follows it, writes a body without that check. Find where, and close the gap.

## Acceptance criteria
- [ ] The path the dashboard's new-item form takes to write a story's or epic's body is traced, and the step that writes without the wip markdown lint is named in the story's notes.
- [ ] Making or editing an item from the dashboard with a body the lint rejects is refused, and the form shows the rule and the line as the CLI does.
- [ ] A test on the side that was missing the check reproduces S-0231's body and fails without the fix.
- [ ] I-0055 is closed with `flai issue close`.

## Tasks
- T-0708 mdlint reports MD007 ul-indent as markdownlint does
- T-0709 The dashboard's new-item and save paths refuse S-0231's body naming MD007 and its line
- T-0710 The design names MD007 among the rules flai lints and I-0055 is closed

## Notes

The trace (T-0708, 2026-10-02). The new-item form posts to `flaiover/src/routes/api/items/+server.ts`, which runs the host action `item.new` (`flai/internal/hostapi/writes.go`), which runs `flai story new --body-stdin`: `itemnew.Create` → `workitem.Repo.Create` → `LintGuard` → `mdlint.Guard`, the same step `flai task new` takes. An edit that follows goes through `doc.save` → `docedit.Save`, which refuses what `flai check`, and so the wip lint, finds new. No step on the dashboard's side writes without the lint.

The step that let S-0231 through is the lint itself: `flai/internal/mdlint` did not implement MD007, unordered list indentation, which is not among the rules ADR-0061 names. Linted with this repository's configuration, S-0231's original body (`git show 4a2ccf0^:wip/kanban/stories/S-0231-author-license-md-file.md`) gave no finding, so `flai story new` and `flai task new` would have taken it too: the goal's premise that they refuse it did not hold. The gap is closed by adding MD007 to mdlint (T-0708), which now reports S-0231's four lines as markdownlint-cli2 does.
