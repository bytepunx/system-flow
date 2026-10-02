---
id: S-0240
type: story
nature: remediation
title: A story made from the dashboard is refused when its body fails the markdown lint, as flai task new refuses it
status: ready
owner: arobson
created: 2026-10-02T12:42:51Z
updated: 2026-10-02T16:24:25Z
transitions:
  - to: ready
    at: 2026-10-02T16:24:25Z
    by: alex
tags: [dashboard, cli]
topics: [dashboard, markdown]
touches: [flaiover/src/lib/server, flai/internal/itemnew, flai/internal/hostapi]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

## Notes
